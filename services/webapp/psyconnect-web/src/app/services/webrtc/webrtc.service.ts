import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';

export interface WebRTCConfig {
  iceServers: RTCIceServer[];
}

export interface CallState {
  isActive: boolean;
  isCalling: boolean;
  isReceivingCall: boolean;
  sessionId?: string;
  remoteUserId?: string;
  localStream?: MediaStream;
  remoteStream?: MediaStream;
}

@Injectable({
  providedIn: 'root',
})
export class WebRTCService {
  private peerConnection: RTCPeerConnection | null = null;
  private localStream: MediaStream | null = null;
  private remoteStream: MediaStream | null = null;
  private pendingIceCandidates: RTCIceCandidateInit[] = []; // Queue for ICE candidates

  private callStateSubject = new BehaviorSubject<CallState>({
    isActive: false,
    isCalling: false,
    isReceivingCall: false,
    localStream: undefined,
    remoteStream: undefined,
  });

  public callState$ = this.callStateSubject.asObservable();

  // Event when call ends (for UI cleanup)
  private callEndedSubject = new Subject<void>();
  public callEnded$ = this.callEndedSubject.asObservable();

  private config: RTCConfiguration = {
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' },
      { urls: 'stun:stun1.l.google.com:19302' },
    ],
  };

  constructor() {
    console.log('📞 WebRTC Service initialized');
  }

  // Initialize local media stream
  async initLocalStream(
    constraints: MediaStreamConstraints = { video: true, audio: true }
  ): Promise<MediaStream> {
    try {
      this.localStream = await navigator.mediaDevices.getUserMedia(constraints);
      console.log('✅ Local stream initialized', this.localStream);

      this.updateCallState({ localStream: this.localStream });
      return this.localStream;
    } catch (error: any) {
      console.error('❌ Error getting user media:', error);

      // If video device not found, try audio only
      if (error.name === 'NotFoundError' && constraints.video) {
        console.log('⚠️ Video device not found, trying audio only...');
        try {
          this.localStream = await navigator.mediaDevices.getUserMedia({
            video: false,
            audio: true,
          });
          console.log('✅ Audio-only stream initialized');
          this.updateCallState({ localStream: this.localStream });
          alert('Camera không khả dụng. Chuyển sang chế độ chỉ có âm thanh.');
          return this.localStream;
        } catch (audioError) {
          console.error('❌ Audio also failed:', audioError);
        }
      }

      throw error;
    }
  }

  // Create peer connection
  createPeerConnection(
    onIceCandidate: (candidate: RTCIceCandidate) => void
  ): RTCPeerConnection {
    this.peerConnection = new RTCPeerConnection(this.config);

    // Handle ICE candidates
    this.peerConnection.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('🧊 ICE candidate generated:', event.candidate);
        onIceCandidate(event.candidate);
      }
    };

    // Handle remote stream
    this.peerConnection.ontrack = (event) => {
      console.log('📺 Remote track received!');
      console.log('📺 Track kind:', event.track.kind);
      console.log('📺 Track label:', event.track.label);
      console.log('📺 Streams:', event.streams);
      console.log('📺 Stream[0]:', event.streams[0]);
      console.log('📺 Stream tracks:', event.streams[0]?.getTracks());

      this.remoteStream = event.streams[0];
      console.log('✅ Remote stream set:', this.remoteStream);

      this.updateCallState({
        remoteStream: this.remoteStream,
        isActive: true,
      });
    };

    // Handle connection state changes
    this.peerConnection.onconnectionstatechange = () => {
      console.log('🔄 Connection state:', this.peerConnection?.connectionState);

      if (this.peerConnection?.connectionState === 'connected') {
        console.log('✅ WebRTC connection established');
        this.updateCallState({ isActive: true, isCalling: false });
      } else if (
        this.peerConnection?.connectionState === 'disconnected' ||
        this.peerConnection?.connectionState === 'failed'
      ) {
        console.log('❌ WebRTC connection failed/disconnected');
        this.endCall();
      }
    };

    // Add local stream tracks
    if (this.localStream) {
      console.log('➕ Adding local tracks to peer connection...');
      const tracks = this.localStream.getTracks();
      console.log(
        `➕ Found ${tracks.length} local tracks:`,
        tracks.map((t) => `${t.kind} (${t.label})`)
      );

      tracks.forEach((track) => {
        console.log(`➕ Adding ${track.kind} track:`, track.label);
        this.peerConnection!.addTrack(track, this.localStream!);
      });
      console.log('✅ All local tracks added to peer connection');
    } else {
      console.warn(
        '⚠️  No local stream available when creating peer connection'
      );
    }

    return this.peerConnection;
  }

  // Create offer (caller)
  async createOffer(): Promise<RTCSessionDescriptionInit> {
    if (!this.peerConnection) {
      throw new Error('Peer connection not initialized');
    }

    const offer = await this.peerConnection.createOffer();
    await this.peerConnection.setLocalDescription(offer);

    console.log('📤 Offer created:', offer);
    this.updateCallState({ isCalling: true });

    return offer;
  }

  // Handle incoming offer (receiver)
  async handleOffer(offer: RTCSessionDescriptionInit): Promise<void> {
    if (!this.peerConnection) {
      throw new Error('Peer connection not initialized');
    }

    await this.peerConnection.setRemoteDescription(
      new RTCSessionDescription(offer)
    );
    console.log('📥 Offer received and set');

    // Process queued ICE candidates
    console.log(
      `🧊 Processing ${this.pendingIceCandidates.length} queued ICE candidates...`
    );
    for (const candidate of this.pendingIceCandidates) {
      try {
        await this.peerConnection.addIceCandidate(
          new RTCIceCandidate(candidate)
        );
        console.log('✅ Queued ICE candidate added');
      } catch (error) {
        console.error('❌ Error adding queued ICE candidate:', error);
      }
    }
    this.pendingIceCandidates = []; // Clear queue

    this.updateCallState({
      isReceivingCall: true,
    });
  }

  // Create answer (callee)
  async createAnswer(): Promise<RTCSessionDescriptionInit> {
    if (!this.peerConnection) {
      throw new Error('Peer connection not initialized');
    }

    const answer = await this.peerConnection.createAnswer();
    await this.peerConnection.setLocalDescription(answer);

    console.log('📤 Answer created:', answer);
    this.updateCallState({ isReceivingCall: false, isCalling: false });

    return answer;
  }

  // Handle incoming answer (caller)
  async handleAnswer(answer: RTCSessionDescriptionInit): Promise<void> {
    if (!this.peerConnection) {
      console.error('❌ Peer connection not initialized');
      return;
    }

    // Check if we're in the right state to receive an answer
    const signalingState = this.peerConnection.signalingState;
    console.log('🔍 Current signaling state:', signalingState);

    if (signalingState !== 'have-local-offer') {
      console.warn('⚠️  Cannot set remote answer in state:', signalingState);
      return; // Ignore duplicate or out-of-order answers
    }

    try {
      await this.peerConnection.setRemoteDescription(
        new RTCSessionDescription(answer)
      );
      console.log('📥 Answer received and set');
      this.updateCallState({ isCalling: false });
    } catch (error) {
      console.error('❌ Error setting remote answer:', error);
    }
  }

  // Add ICE candidate
  async addIceCandidate(candidate: RTCIceCandidateInit): Promise<void> {
    if (!this.peerConnection) {
      console.warn('⚠️  Peer connection not ready for ICE candidate');
      return;
    }

    // Check if remote description is set
    if (!this.peerConnection.remoteDescription) {
      console.log('🧊 Queuing ICE candidate (remote description not set yet)');
      this.pendingIceCandidates.push(candidate);
      return;
    }

    try {
      await this.peerConnection.addIceCandidate(new RTCIceCandidate(candidate));
      console.log('✅ ICE candidate added');
    } catch (error) {
      console.error('❌ Error adding ICE candidate:', error);
    }
  }

  // Toggle video
  toggleVideo(): boolean {
    if (!this.localStream) return false;

    const videoTrack = this.localStream.getVideoTracks()[0];
    if (videoTrack) {
      videoTrack.enabled = !videoTrack.enabled;
      console.log('📹 Video:', videoTrack.enabled ? 'ON' : 'OFF');
      return videoTrack.enabled;
    }
    return false;
  }

  // Toggle audio
  toggleAudio(): boolean {
    if (!this.localStream) return false;

    const audioTrack = this.localStream.getAudioTracks()[0];
    if (audioTrack) {
      audioTrack.enabled = !audioTrack.enabled;
      console.log('🎤 Audio:', audioTrack.enabled ? 'ON' : 'OFF');
      return audioTrack.enabled;
    }
    return false;
  }

  // End call
  endCall(): void {
    console.log('📞 Ending call...');

    // Stop local stream
    if (this.localStream) {
      this.localStream.getTracks().forEach((track) => {
        track.stop();
        console.log('🛑 Stopped track:', track.kind);
      });
      this.localStream = null;
    }

    // Close peer connection
    if (this.peerConnection) {
      this.peerConnection.close();
      this.peerConnection = null;
    }

    // Reset remote stream
    this.remoteStream = null;

    // Clear pending ICE candidates
    this.pendingIceCandidates = [];

    // Reset state
    this.callStateSubject.next({
      isActive: false,
      isCalling: false,
      isReceivingCall: false,
      sessionId: undefined,
      remoteUserId: undefined,
      localStream: undefined,
      remoteStream: undefined,
    });

    // Emit call ended event for UI cleanup
    this.callEndedSubject.next();

    console.log('✅ Call ended');
  }

  // Get current call state
  getCallState(): CallState {
    return this.callStateSubject.value;
  }

  // Update call state
  private updateCallState(updates: Partial<CallState>): void {
    const currentState = this.callStateSubject.value;
    this.callStateSubject.next({ ...currentState, ...updates });
  }

  // Set session info
  setSessionInfo(sessionId: string, remoteUserId: string): void {
    this.updateCallState({ sessionId, remoteUserId });
  }

  // Get local stream
  getLocalStream(): MediaStream | null {
    return this.localStream;
  }

  // Get remote stream
  getRemoteStream(): MediaStream | null {
    return this.remoteStream;
  }

  // Get available media devices
  async getAvailableDevices(): Promise<{
    cameras: MediaDeviceInfo[];
    microphones: MediaDeviceInfo[];
  }> {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices();
      return {
        cameras: devices.filter((d) => d.kind === 'videoinput'),
        microphones: devices.filter((d) => d.kind === 'audioinput'),
      };
    } catch (error) {
      console.error('❌ Error getting devices:', error);
      return { cameras: [], microphones: [] };
    }
  }

  // Switch camera
  async switchCamera(deviceId: string): Promise<boolean> {
    if (!this.localStream) return false;

    try {
      const newStream = await navigator.mediaDevices.getUserMedia({
        video: { deviceId: { exact: deviceId } },
        audio: false,
      });

      const newVideoTrack = newStream.getVideoTracks()[0];
      const oldVideoTrack = this.localStream.getVideoTracks()[0];

      if (this.peerConnection) {
        const sender = this.peerConnection
          .getSenders()
          .find((s) => s.track?.kind === 'video');
        if (sender) await sender.replaceTrack(newVideoTrack);
      }

      this.localStream.removeTrack(oldVideoTrack);
      this.localStream.addTrack(newVideoTrack);
      oldVideoTrack.stop();

      this.updateCallState({ localStream: this.localStream });
      console.log('✅ Camera switched');
      return true;
    } catch (error) {
      console.error('❌ Error switching camera:', error);
      return false;
    }
  }

  // Switch microphone
  async switchMicrophone(deviceId: string): Promise<boolean> {
    if (!this.localStream) return false;

    try {
      const newStream = await navigator.mediaDevices.getUserMedia({
        video: false,
        audio: { deviceId: { exact: deviceId } },
      });

      const newAudioTrack = newStream.getAudioTracks()[0];
      const oldAudioTrack = this.localStream.getAudioTracks()[0];

      if (this.peerConnection) {
        const sender = this.peerConnection
          .getSenders()
          .find((s) => s.track?.kind === 'audio');
        if (sender) await sender.replaceTrack(newAudioTrack);
      }

      this.localStream.removeTrack(oldAudioTrack);
      this.localStream.addTrack(newAudioTrack);
      oldAudioTrack.stop();

      this.updateCallState({ localStream: this.localStream });
      console.log('✅ Microphone switched');
      return true;
    } catch (error) {
      console.error('❌ Error switching microphone:', error);
      return false;
    }
  }

  // Get current device IDs
  getCurrentDevices(): { cameraId?: string; microphoneId?: string } {
    if (!this.localStream) return {};

    const videoTrack = this.localStream.getVideoTracks()[0];
    const audioTrack = this.localStream.getAudioTracks()[0];

    return {
      cameraId: videoTrack?.getSettings().deviceId,
      microphoneId: audioTrack?.getSettings().deviceId,
    };
  }
}
