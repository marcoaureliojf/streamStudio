import axios from 'axios';
import React, { useEffect, useRef, useState } from 'react';
import { useAuth } from './Auth';

const Stream = () => {
     const { token } = useAuth();
    const localVideo = useRef();
 const [localStream, setLocalStream] = useState(null);
     const [remoteSDP, setRemoteSDP] = useState(null);

     useEffect(() => {
       const startStream = async () => {
         try {
            const stream = await navigator.mediaDevices.getUserMedia({video: true, audio: true});
             localVideo.current.srcObject = stream;
           setLocalStream(stream);
           } catch (err) {
                console.error("Erro ao iniciar stream", err)
            }
       }
           startStream()
       },[])

   const handleOffer = async () => {
    if (!localStream) {
        console.log('Local stream not available');
        return;
    }
     
    const peerConnection = new RTCPeerConnection({
      iceServers: [{urls: 'stun:stun.l.google.com:19302'}]
    })
    
       for (const track of localStream.getTracks()) {
            peerConnection.addTrack(track, localStream);
       }
     
      peerConnection.onicecandidate = (event) => {
        if (event.candidate) {
           console.log("ICE candidate:", event.candidate);
         }
      };
    
       peerConnection.ontrack = event => {
           console.log("Track", event)
     }
    
     
    const offer = await peerConnection.createOffer();
  
    await peerConnection.setLocalDescription(offer)
   
    try {
      const response = await axios.post(
          'http://localhost:8081/api/offer',
          {
            sdp: offer.sdp,
          },
          {
            headers: {
                Authorization: `Bearer ${token}`,
           },
        }
       );
     setRemoteSDP(response.data.sdp)
       await peerConnection.setRemoteDescription(
           new RTCSessionDescription({
                type: 'answer',
                sdp: response.data.sdp,
            })
       );
    } catch (error) {
       console.error('Error ao enviar o offer', error);
    }
   }

  return (
    <div>
      <video ref={localVideo} autoPlay muted width="320" height="240" />
       <button type="button" onClick={handleOffer}>Iniciar Transmissão</button>
    </div>
  );
};

export default Stream;