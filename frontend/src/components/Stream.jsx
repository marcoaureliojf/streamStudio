// import axios from 'axios';
// import React, { useEffect, useRef, useState } from 'react';
// import { useAuth } from './Auth';

// const Stream = () => {
//     const { token } = useAuth();
//     const localVideo = useRef();
//     const remoteVideo = useRef();
//     const [peerConnection, setPeerConnection] = useState(null);
//     const [localStream, setLocalStream] = useState(null);
//     const [videoDevices, setVideoDevices] = useState([]);
//     const [audioDevices, setAudioDevices] = useState([]);
//     const [selectedVideoDevice, setSelectedVideoDevice] = useState('');
//     const [selectedAudioDevice, setSelectedAudioDevice] = useState('');
//     const [permissionError, setPermissionError] = useState('');

//     useEffect(() => {
//         // const getMediaDevices = async () => {
//         //     try {
//         //         // Request permissions
//         //         await navigator.mediaDevices.getUserMedia({ video: true, audio: true });

//         //         // Enumerate devices
//         //         const devices = await navigator.mediaDevices.enumerateDevices();
//         //         const videoDevices = devices.filter(device => device.kind === 'videoinput');
//         //         const audioDevices = devices.filter(device => device.kind === 'audioinput');
//         //         setVideoDevices(videoDevices);
//         //         setAudioDevices(audioDevices);
//         //         if (videoDevices.length > 0) setSelectedVideoDevice(videoDevices[0].deviceId);
//         //         if (audioDevices.length > 0) setSelectedAudioDevice(audioDevices[0].deviceId);
//         //         console.log("Video Devices:", videoDevices);
//         //         console.log("Audio Devices:", audioDevices);
//         //     } catch (err) {
//         //         console.error("Error getting media devices", err);
//         //         setPermissionError('Permission denied. Please allow access to camera and microphone.');
//         //     }
//         // };

//         // getMediaDevices();
//         const getMediaDevices = async () => {
//             try {
//                 await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
//                 const devices = await navigator.mediaDevices.enumerateDevices();
//                 const video = devices.filter(device => device.kind === 'videoinput')
//                 const audio = devices.filter(device => device.kind === 'audioinput')
//                 setVideoDevices(video);
//                setAudioDevices(audio);
//                if (video.length > 0) {
//                    setSelectedVideoDevice(video[0].deviceId);
//                 }
//                 if(audio.length > 0){
//                   setSelectedAudioDevice(audio[0].deviceId);
//                }
//              } catch (err) {
//                    console.error("Erro ao enumerar dispositivos de mídia", err)
//               }
//         };
//        getMediaDevices();
//     }, []);

//     useEffect(() => {
//         const startStream = async () => {
//             if (selectedVideoDevice && selectedAudioDevice) {
//                 try {
//                     const stream = await navigator.mediaDevices.getUserMedia({
//                         video: { deviceId: selectedVideoDevice },
//                         audio: { deviceId: selectedAudioDevice }
//                     });
//                     localVideo.current.srcObject = stream;
//                     setLocalStream(stream);
//                 } catch (err) {
//                     console.error("Error starting stream", err);
//                 }
//             }
//         };

//         startStream();
//     }, [selectedVideoDevice, selectedAudioDevice]);

//     // const startStream = async () => {
//     //     if (selectedVideoDevice && selectedAudioDevice) {
//     //         try {
//     //             const stream = await navigator.mediaDevices.getUserMedia({
//     //                 video: { deviceId: selectedVideoDevice },
//     //                 audio: { deviceId: selectedAudioDevice }
//     //             });
//     //             localVideo.current.srcObject = stream;
//     //             setLocalStream(stream);
//     //         } catch (err) {
//     //             console.error("Erro ao iniciar stream:", err);
//     //             if (err.name === 'NotAllowedError') {
//     //               alert("Permissão negada para acessar a câmera e o microfone. Por favor, verifique as permissões do seu navegador.");
//     //           }
//     //         }
//     //     }
//     // };

//     // const stopStream = () => {
//     //     if (localStream) {
//     //         localStream.getTracks().forEach(track => track.stop());
//     //         localVideo.current.srcObject = null;
//     //         setLocalStream(null);
//     //     }
//     // };

//     const handleVideoDeviceChange = (event) => {
//         setSelectedVideoDevice(event.target.value);
//     };

//     const handleAudioDeviceChange = (event) => {
//         setSelectedAudioDevice(event.target.value);
//     };

//     const handleOffer = async () => {
//         if (!localStream) {
//            console.log('Local stream not available');
//             return;
//         }

//         const pc = new RTCPeerConnection({
//             iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
//        });

//         setPeerConnection(pc);

//        localStream.getTracks().forEach(track => {
//          pc.addTrack(track, localStream);
//         });

//         pc.onicecandidate = async (event) => {
//             if (event.candidate) {
//                 console.log("ICE candidate:", event.candidate);
//                try {
//                     await axios.post(
//                         'http://localhost:8182/api/icecandidate',
//                      {
//                           candidate: event.candidate,
//                        },
//                         {
//                          headers: {
//                               Authorization: `Bearer ${token}`,
//                            },
//                         }
//                    );
//                }
//                 catch (err) {
//                    console.log("Erro ao enviar candidato ICE:", err);
//                }
//             }
//         };
    
//        pc.ontrack =  event => {
//             console.log("Track:", event.track)
//             if (remoteVideo.current) {
//                  remoteVideo.current.srcObject = event.streams[0];
//             }
//      };


    
//     try {
//         const offer = await pc.createOffer();
//         await pc.setLocalDescription(offer);

//         const response = await axios.post(
//              'http://localhost:8182/api/offer',
//              {
//                    sdp: offer.sdp,
//                 },
//             {
//                     headers: {
//                         Authorization: `Bearer ${token}`,
//                     },
//                }
//            );
    
//         // const answer = new RTCSessionDescription({
//         //       type: 'answer',
//         //         sdp: response.data.sdp,
//         //    });
//         //    await pc.setRemoteDescription(answer);
//         const { sdp } = response.data;
//         const remoteDesc = new RTCSessionDescription({ type: 'answer', sdp });
//         await pc.setRemoteDescription(remoteDesc);  
//         } catch (error) {
//           console.error('Error ao enviar o offer', error);
//       }
//     };

//     const handleDisconnect = () => {
//         if (peerConnection) {
//            peerConnection.close();
//            console.log("Desconectado")
//        }
//   };

//     // return (
//     //     <div>
//     //         {permissionError && <p style={{ color: 'red' }}>{permissionError}</p>}
//     //         <div>
//     //             <label htmlFor="videoDevices">Select Video Device:</label>
//     //             <select id="videoDevices" onChange={handleVideoDeviceChange} value={selectedVideoDevice || ''}>
//     //                 {videoDevices.map(device => (
//     //                     <option key={device.deviceId} value={device.deviceId}>{device.label || `Camera ${device.deviceId}`}</option>
//     //                 ))}
//     //             </select>
//     //         </div>
//     //         <div>
//     //             <label htmlFor="audioDevices">Select Audio Device:</label>
//     //             <select id="audioDevices" onChange={handleAudioDeviceChange} value={selectedAudioDevice || ''}>
//     //                 {audioDevices.map(device => (
//     //                     <option key={device.deviceId} value={device.deviceId}>{device.label || `Microphone ${device.deviceId}`}</option>
//     //                 ))}
//     //             </select>
//     //         </div>
//     //         <div>
//     //             <button onClick={handleOffer}>Start Stream</button>
//     //             <button onClick={handleDisconnect}>Stop Stream</button>
//     //         </div>
//     //         <video ref={localVideo} autoPlay muted width="320" height="240" />
//     //         <video ref={remoteVideo} autoPlay width="320" height="240">
//     //             <track kind="captions" src="" label="English captions" srcLang="en" />
//     //         </video>
//     //     </div>
//     // );
//     return (
//         <div>
//               <div>
//                <label htmlFor="videoDevice">Câmera:</label>
//                    <select id="videoDevice" onChange={handleVideoDeviceChange} value={selectedVideoDevice}>
//                        {videoDevices.map((device) => (
//                           <option key={device.deviceId} value={device.deviceId}>
//                                {device.label}
//                             </option>
//                        ))}
//                   </select>
//                </div>
//                <div>
//                   <label htmlFor="audioDevice">Microfone:</label>
//                   <select id="audioDevice" onChange={handleAudioDeviceChange} value={selectedAudioDevice}>
//                     {audioDevices.map((device) => (
//                         <option key={device.deviceId} value={device.deviceId}>
//                               {device.label}
//                           </option>
//                       ))}
//                </select>
//            </div>
//         <video ref={localVideo} autoPlay muted width="320" height="240" />
//         <video ref={remoteVideo} autoPlay  width="320" height="240" />
//         <button onClick={handleOffer}>Iniciar Transmissão</button>
//           <button onClick={handleDisconnect}>Desconectar</button>
//       </div>
//   );
// };

// export default Stream;

import axios from 'axios';
import React, { useEffect, useRef, useState } from 'react';
import { useAuth } from './Auth';

const Stream = () => {
    const { token } = useAuth();
    const localVideo = useRef(null);
     const remoteVideo = useRef(null);
    const [localStream, setLocalStream] = useState(null);
    const [videoDevices, setVideoDevices] = useState([]);
   const [audioDevices, setAudioDevices] = useState([]);
   const [selectedVideoDevice, setSelectedVideoDevice] = useState('');
     const [selectedAudioDevice, setSelectedAudioDevice] = useState('');
     const pc = useRef(null);

    useEffect(() => {
      const getMediaDevices = async () => {
           try {
                await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
                const devices = await navigator.mediaDevices.enumerateDevices();
                const video = devices.filter(device => device.kind === 'videoinput')
                const audio = devices.filter(device => device.kind === 'audioinput')
                   setVideoDevices(video);
                   setAudioDevices(audio);
                    if (video.length > 0) {
                      setSelectedVideoDevice(video[0].deviceId);
                    }
                 if (audio.length > 0) {
                     setSelectedAudioDevice(audio[0].deviceId);
                   }
           } catch (err) {
                console.error("Erro ao enumerar dispositivos de mídia", err)
              }
       }
         getMediaDevices()
   }, []);
      
    useEffect(() => {
    const startStream = async () => {
        if (!selectedVideoDevice || !selectedAudioDevice) {
            return;
        }
    
     try {
            const stream = await navigator.mediaDevices.getUserMedia({
                 video: { deviceId: selectedVideoDevice },
                audio: { deviceId: selectedAudioDevice },
            });
           localVideo.current.srcObject = stream;
             setLocalStream(stream);
        } catch (err) {
              console.error("Erro ao iniciar stream:", err);
             if (err.name === 'NotAllowedError') {
                 alert("Permissão negada para acessar a câmera e o microfone. Por favor, verifique as permissões do seu navegador.");
             }
         }
     };
    
       startStream();
     }, [selectedVideoDevice, selectedAudioDevice]);
        
      const handleVideoDeviceChange = (e) => {
           setSelectedVideoDevice(e.target.value);
       }
    
      const handleAudioDeviceChange = (e) => {
           setSelectedAudioDevice(e.target.value);
      }
        
      const handleOffer = async () => {
          if (!localStream) {
              console.log('Local stream not available');
                return;
           }
           
            pc.current = new RTCPeerConnection({
                 iceServers: [{ urls: 'stun:stun.l.google.com:19302' }]
            });
       
        localStream.getTracks().forEach(track => {
               pc.current.addTrack(track, localStream);
        });

        pc.current.onicecandidate = async (event) => {
             if (event.candidate) {
                   console.log("ICE candidate:", event.candidate);
                   try {
                        await axios.post(
                            'http://localhost:8182/api/icecandidate',
                            {
                                candidate: event.candidate,
                            },
                            {
                                headers: {
                                   Authorization: `Bearer ${token}`,
                                 },
                            }
                         );
                    } catch (err) {
                       console.log("Erro ao enviar candidato ICE:", err);
                 }
             }
          };
    
         pc.current.ontrack = (event) => {
             if (remoteVideo.current && event.streams && event.streams[0]) {
                  remoteVideo.current.srcObject = event.streams[0];
              }
         };
    
          try {
            const offer = await pc.current.createOffer();
             await pc.current.setLocalDescription(offer);
             
             const response = await axios.post(
                   'http://localhost:8182/api/offer',
                    {
                      sdp: offer.sdp,
                     },
                     {
                        headers: {
                            Authorization: `Bearer ${token}`,
                          },
                     }
                );
        
              const answer = new RTCSessionDescription({
                   type: 'answer',
                    sdp: response.data.sdp,
                });
             await pc.current.setRemoteDescription(answer)
          } catch (err) {
            console.error("Erro ao iniciar a transmissão:", err)
         }
        };
            
         const handleDisconnect = () => {
           if(pc.current){
                pc.current.close();
               console.log("Desconectado")
             }
        };
            
            return (
               <div>
                  <div>
                        <label htmlFor="videoDevice">Câmera:</label>
                        <select id="videoDevice" onChange={handleVideoDeviceChange} value={selectedVideoDevice}>
                            {videoDevices.map((device) => (
                               <option key={device.deviceId} value={device.deviceId}>
                                   {device.label}
                                </option>
                            ))}
                         </select>
                    </div>
                   <div>
                        <label htmlFor="audioDevice">Microfone:</label>
                      <select id="audioDevice" onChange={handleAudioDeviceChange} value={selectedAudioDevice}>
                           {audioDevices.map((device) => (
                              <option key={device.deviceId} value={device.deviceId}>
                                   {device.label}
                               </option>
                            ))}
                       </select>
                   </div>
                 <video ref={localVideo} autoPlay muted width="320" height="240" />
                  <video ref={remoteVideo} autoPlay width="320" height="240" />
                  <button onClick={handleOffer}>Iniciar Transmissão</button>
                    <button onClick={handleDisconnect}>Desconectar</button>
               </div>
            );
     };

    export default Stream;