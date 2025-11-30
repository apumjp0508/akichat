import { sendOffer, sendIce, sendAnswer } from "./signaling";
import { getLocalStream } from "./callSession";

export async function startCall(
  ws: WebSocket,
  toUserId: number,
  localStream?: MediaStream
): Promise<RTCPeerConnection> {
  const pc = new RTCPeerConnection({
    //ここで通信経路を探し出してくれるサーバーを設定している。
    iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
  });

  //localstream はカメラやマイクのこと

  if (localStream) {
    for (const track of localStream.getTracks()) {
      pc.addTrack(track, localStream);
    }
  }

  //onicecandidate イベントは、その候補経路を見つけるたびに発火し、

  pc.onicecandidate = (ev) => {
    console.log("ev:", ev);
    if (ev.candidate) {
      const cand: RTCIceCandidateInit = {
        candidate: ev.candidate.candidate,
        sdpMid: ev.candidate.sdpMid ?? undefined,
        sdpMLineIndex: ev.candidate.sdpMLineIndex ?? undefined,
      };
      sendIce(ws, toUserId, cand);
    }
  };

  const offer = await pc.createOffer();
  await pc.setLocalDescription(offer);
  sendOffer(ws, toUserId, offer);

  return pc;
}

// 受信側: オファーを受け取り、アンサーを返す
export async function acceptCall(
  ws: WebSocket,
  fromUserId: number,
  offerSdp: RTCSessionDescriptionInit,
): Promise<RTCPeerConnection> {
  const pc = new RTCPeerConnection({
    iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
  });

  const localStream = await getLocalStream();

  if (localStream) {
    for (const track of localStream.getTracks()) {
      pc.addTrack(track, localStream);
    }
  }

  pc.onicecandidate = (ev) => {
    if (ev.candidate) {
      const cand: RTCIceCandidateInit = {
        candidate: ev.candidate.candidate,
        sdpMid: ev.candidate.sdpMid ?? undefined,
        sdpMLineIndex: ev.candidate.sdpMLineIndex ?? undefined,
      };
      sendIce(ws, fromUserId, cand);
    }
  };

  await pc.setRemoteDescription(offerSdp);
  const answer = await pc.createAnswer();
  await pc.setLocalDescription(answer);
  sendAnswer(ws, fromUserId, answer);

  return pc;
}

// （ICE適用は callSession.ts の applyRemoteIce に統一）

