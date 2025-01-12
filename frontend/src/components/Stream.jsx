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

    // Para controlar se já obteve permissões e dispositivos
    const [permissionsObtained, setPermissionsObtained] = useState(false);

    useEffect(() => {
        // Apenas chama uma vez
        const getMediaDevices = async () => {
            if (permissionsObtained) return; // Evita chamadas duplicadas

            try {
                // Solicita permissão para acessar mídia
                await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
                const devices = await navigator.mediaDevices.enumerateDevices();
                console.log('Dispositivos encontrados:', devices);

                const video = devices.filter(device => device.kind === 'videoinput');
                const audio = devices.filter(device => device.kind === 'audioinput');

                setVideoDevices(video);
                setAudioDevices(audio);

                if (video.length > 0) {
                    setSelectedVideoDevice(video[0].deviceId);
                }
                if (audio.length > 0) {
                    setSelectedAudioDevice(audio[0].deviceId);
                }

                setPermissionsObtained(true); // Marca como permissão obtida
            } catch (err) {
                console.error("Erro ao enumerar dispositivos de mídia", err);
            }
        };

        getMediaDevices(); // Chama apenas uma vez
    }, [permissionsObtained]);

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
                    alert("Permissão negada para acessar a câmera e o microfone.");
                }
            }
        };

        startStream();
    }, [selectedVideoDevice, selectedAudioDevice]);

    const handleVideoDeviceChange = (e) => {
        setSelectedVideoDevice(e.target.value);
    };

    const handleAudioDeviceChange = (e) => {
        setSelectedAudioDevice(e.target.value);
    };

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
            console.log('Track recebido:', event.track);
            console.log('Stream associado:', event.streams[0]);

            if (remoteVideo.current && event.streams && event.streams[0]) {
                console.log('Associando stream remoto ao vídeo');
                remoteVideo.current.srcObject = event.streams[0];
            } else {
                console.error('Stream remoto ou referência de vídeo não encontrada.');
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
            await pc.current.setRemoteDescription(answer);
        } catch (err) {
            console.error("Erro ao iniciar a transmissão:", err);
        }
    };

    const handleDisconnect = () => {
        if (pc.current) {
            pc.current.close();
            console.log("Desconectado");
        }
    };

    return (
        <div>
            <div>
                <label htmlFor="videoDevice">Câmera:</label>
                <select id="videoDevice" onChange={handleVideoDeviceChange} value={selectedVideoDevice}>
                    {videoDevices.map((device) => (
                        <option key={device.deviceId} value={device.deviceId}>
                            {device.label || 'Sem nome'}
                        </option>
                    ))}
                </select>
            </div>
            <div>
                <label htmlFor="audioDevice">Microfone:</label>
                <select id="audioDevice" onChange={handleAudioDeviceChange} value={selectedAudioDevice}>
                    {audioDevices.map((device) => (
                        <option key={device.deviceId} value={device.deviceId}>
                            {device.label || 'Sem nome'}
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
