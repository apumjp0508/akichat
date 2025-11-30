import { startCall } from "./startChatOffer";

type ConnState = RTCPeerConnectionState | "none";

let currentPc: RTCPeerConnection | null = null;
let currentPeerUserId: number | null = null;
const listeners = new Set<(s: ConnState) => void>();

function notify() {
  const state: ConnState = currentPc ? currentPc.connectionState : "none";
  for (const l of listeners) l(state);
}

export async function getLocalStream(): Promise<MediaStream> {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: true,
      audio: true,
    });
    return stream;
  } catch (e) {
    console.error("getUserMedia error:", e);
    throw e;
  }
}


export async function startOutgoingCall(
  ws: WebSocket,
  toUserId: number,
) {
  const localStream = await getLocalStream();
  currentPc = await startCall(ws, toUserId, localStream);
  currentPeerUserId = toUserId;
  currentPc.onconnectionstatechange = notify;
  notify();
  return currentPc;
}

//受け取った Answer が今の相手か検査
export async function handleRemoteAnswer(
  fromUserId: number,
  sdp: RTCSessionDescriptionInit
) {
  if (!currentPc || currentPeerUserId !== fromUserId) return;
  if (currentPc.signalingState === "stable") return;
  await currentPc.setRemoteDescription(sdp);
  notify();
}

// 統一: 受信/送信どちらでも、fromUserId と状況に応じて ICE を適用
export async function applyRemoteIce(
  fromUserId: number,
  candidate: RTCIceCandidateInit,
  peerByUser?: Map<number, RTCPeerConnection>
) {
  let pc: RTCPeerConnection | null | undefined = null;
  if (peerByUser) {
    pc = peerByUser.get(fromUserId) ?? null;
  } else if (currentPeerUserId === fromUserId) {
    pc = currentPc;
  }
  if (!pc) return;
  await pc.addIceCandidate(candidate);
}

// 互換: 既存の呼び出しをサポート（内部で統一関数を利用）
export async function applyRemoteIceForOutgoing(
  fromUserId: number,
  candidate: RTCIceCandidateInit
) {
  await applyRemoteIce(fromUserId, candidate);
}

export function subscribeConnectionState(
  cb: (state: ConnState) => void
): () => void {
  listeners.add(cb);
  cb(currentPc ? currentPc.connectionState : "none");
  return () => listeners.delete(cb);
}

export function getConnectionState(): ConnState {
  return currentPc ? currentPc.connectionState : "none";
}



