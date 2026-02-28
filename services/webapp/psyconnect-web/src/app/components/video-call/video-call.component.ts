import { CommonModule } from '@angular/common';
import {
  AfterViewInit,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import { Subscription } from 'rxjs';
import { CallState, WebRTCService } from '../../services/webrtc/webrtc.service';
import { AvatarFallbackPipe } from '../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'app-video-call',
  standalone: true,
  imports: [CommonModule, AvatarFallbackPipe],
  templateUrl: './video-call.component.html',
  styleUrls: ['./video-call.component.scss'],
})
export class VideoCallComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('localVideo') localVideo!: ElementRef<HTMLVideoElement>;
  @ViewChild('remoteVideo') remoteVideo!: ElementRef<HTMLVideoElement>;

  @Input() conversationId!: string;
  @Input() remoteUserId!: string;
  @Input() currentUserId!: string;
  @Input() remoteUserAvatar?: string;
  @Input() remoteUserName?: string;

  @Output() callEnded = new EventEmitter<void>();
  @Output() iceCandidate = new EventEmitter<RTCIceCandidate>();
  @Output() offer = new EventEmitter<RTCSessionDescriptionInit>();
  @Output() answer = new EventEmitter<RTCSessionDescriptionInit>();

  callState: CallState = {
    isActive: false,
    isCalling: false,
    isReceivingCall: false,
  };

  isVideoEnabled = true;
  isAudioEnabled = true;
  isFullscreen = false;
  isLocalVideoMinimized = false;
  isPipMode = false;
  showSettings = false;

  // Device selection
  availableCameras: MediaDeviceInfo[] = [];
  availableMicrophones: MediaDeviceInfo[] = [];
  selectedCameraId?: string;
  selectedMicrophoneId?: string;

  localVideoPosition = { x: 20, y: 20 };
  pipPosition = { x: window.innerWidth - 340, y: window.innerHeight - 260 };
  private isDragging = false;
  private isDraggingPip = false;
  private dragOffset = { x: 0, y: 0 };

  private callStateSubscription?: Subscription;

  constructor(public webrtcService: WebRTCService) {}

  ngOnInit() {
    console.log('🎥 VideoCallComponent initialized');

    // Add beforeunload warning to prevent accidental tab close during call
    window.addEventListener('beforeunload', this.beforeUnloadHandler);
  }

  ngOnDestroy() {
    this.callStateSubscription?.unsubscribe();
    this.endCall();

    // Remove beforeunload warning
    window.removeEventListener('beforeunload', this.beforeUnloadHandler);
  }

  private beforeUnloadHandler = (e: BeforeUnloadEvent) => {
    // Only show warning if call is active
    if (this.callState.isActive) {
      e.preventDefault();
      e.returnValue = '';
      return 'Bạn đang trong cuộc gọi. Bạn có chắc muốn thoát?';
    }
    return undefined;
  };

  ngAfterViewInit() {
    console.log('🎥 VideoCallComponent view initialized');
    console.log('🎥 Local video element:', this.localVideo?.nativeElement);
    console.log('🎥 Remote video element:', this.remoteVideo?.nativeElement);

    // Load available devices
    this.loadDevices();

    // Subscribe to call state changes
    this.callStateSubscription = this.webrtcService.callState$.subscribe(
      (state) => {
        console.log('📊 Call state updated:', state);
        this.callState = state;

        // Attach streams to video elements
        if (state.localStream && this.localVideo) {
          console.log('🎥 Attaching local stream to video element');
          this.localVideo.nativeElement.srcObject = state.localStream;
        }

        if (state.remoteStream && this.remoteVideo) {
          console.log('🎥 Attaching remote stream to video element');
          this.remoteVideo.nativeElement.srcObject = state.remoteStream;
        }
      }
    );
  }

  // Start outgoing call
  async startCall() {
    try {
      console.log('📞 Starting call to:', this.remoteUserId);

      // Initialize local stream
      await this.webrtcService.initLocalStream();

      // Create peer connection
      this.webrtcService.createPeerConnection((candidate) => {
        this.iceCandidate.emit(candidate);
      });

      // Create and send offer
      const offer = await this.webrtcService.createOffer();
      console.log('📤 Emitting offer event:', offer);
      this.offer.emit(offer);
      console.log('✅ Offer event emitted');
    } catch (error) {
      console.error('❌ Error starting call:', error);
      alert(
        'Không thể bắt đầu cuộc gọi. Vui lòng kiểm tra quyền camera/microphone.'
      );
    }
  }

  // Accept incoming call
  async acceptCall(offer: RTCSessionDescriptionInit) {
    try {
      console.log('📞 Accepting call from:', this.remoteUserId);
      console.log('📞 Offer received:', offer);

      // Initialize local stream
      console.log('🎥 Initializing local stream...');
      await this.webrtcService.initLocalStream();
      console.log('✅ Local stream initialized');

      // Create peer connection
      console.log('🔗 Creating peer connection...');
      this.webrtcService.createPeerConnection((candidate) => {
        console.log('🧊 ICE candidate generated (receiver)');
        this.iceCandidate.emit(candidate);
      });
      console.log('✅ Peer connection created');

      // Handle offer
      console.log('📥 Handling offer...');
      await this.webrtcService.handleOffer(offer);
      console.log('✅ Offer handled');

      // Create and send answer
      console.log('📤 Creating answer...');
      const answer = await this.webrtcService.createAnswer();
      console.log('📤 Answer created:', answer);
      console.log('📤 Emitting answer event');
      this.answer.emit(answer);
      console.log('✅ Answer emitted');
    } catch (error) {
      console.error('❌ Error accepting call:', error);
      alert('Không thể chấp nhận cuộc gọi.');
    }
  }

  // Handle incoming answer
  async handleAnswer(answer: RTCSessionDescriptionInit) {
    await this.webrtcService.handleAnswer(answer);
  }

  // Handle incoming ICE candidate
  async handleIceCandidate(candidate: RTCIceCandidateInit) {
    await this.webrtcService.addIceCandidate(candidate);
  }

  // Toggle video
  toggleVideo() {
    this.isVideoEnabled = this.webrtcService.toggleVideo();
  }

  // Toggle audio
  toggleAudio() {
    this.isAudioEnabled = this.webrtcService.toggleAudio();
  }

  // Toggle fullscreen
  toggleFullscreen() {
    this.isFullscreen = !this.isFullscreen;
  }

  // End call
  endCall() {
    console.log('🔴 End call button clicked');

    // End call in service (stops streams, closes peer connection)
    this.webrtcService.endCall();

    // Emit event to parent to handle UI cleanup and send leave message
    this.callEnded.emit();

    console.log('✅ End call event emitted');
  }

  // Drag functionality for local video
  startDrag(event: MouseEvent) {
    if ((event.target as HTMLElement).closest('.minimize-btn')) {
      return;
    }

    this.isDragging = true;
    this.dragOffset = {
      x: event.clientX - this.localVideoPosition.x,
      y: event.clientY - this.localVideoPosition.y,
    };

    document.addEventListener('mousemove', this.onDrag);
    document.addEventListener('mouseup', this.stopDrag);
    event.preventDefault();
  }

  private onDrag = (event: MouseEvent) => {
    if (!this.isDragging) return;
    this.localVideoPosition = {
      x: event.clientX - this.dragOffset.x,
      y: event.clientY - this.dragOffset.y,
    };
  };

  private stopDrag = () => {
    this.isDragging = false;
    document.removeEventListener('mousemove', this.onDrag);
    document.removeEventListener('mouseup', this.stopDrag);
  };

  toggleLocalVideoSize(event: MouseEvent) {
    event.stopPropagation();
    this.isLocalVideoMinimized = !this.isLocalVideoMinimized;
  }

  // Toggle Picture-in-Picture mode
  togglePipMode() {
    this.isPipMode = !this.isPipMode;
  }

  // Drag functionality for PiP window
  startPipDrag(event: MouseEvent) {
    // Only drag if in PiP mode and not clicking on buttons
    if (!this.isPipMode || (event.target as HTMLElement).closest('button')) {
      return;
    }

    this.isDraggingPip = true;
    this.dragOffset = {
      x: event.clientX - this.pipPosition.x,
      y: event.clientY - this.pipPosition.y,
    };

    document.addEventListener('mousemove', this.onPipDrag);
    document.addEventListener('mouseup', this.stopPipDrag);
    event.preventDefault();
  }

  private onPipDrag = (event: MouseEvent) => {
    if (!this.isDraggingPip) return;

    this.pipPosition = {
      x: event.clientX - this.dragOffset.x,
      y: event.clientY - this.dragOffset.y,
    };
  };

  private stopPipDrag = () => {
    this.isDraggingPip = false;
    document.removeEventListener('mousemove', this.onPipDrag);
    document.removeEventListener('mouseup', this.stopPipDrag);
  };

  // Load available devices
  async loadDevices() {
    const devices = await this.webrtcService.getAvailableDevices();
    this.availableCameras = devices.cameras;
    this.availableMicrophones = devices.microphones;

    // Set current devices
    const current = this.webrtcService.getCurrentDevices();
    this.selectedCameraId = current.cameraId;
    this.selectedMicrophoneId = current.microphoneId;

    console.log('📹 Available cameras:', this.availableCameras.length);
    console.log('🎤 Available microphones:', this.availableMicrophones.length);
  }

  // Camera change handler
  async onCameraChange(deviceId: string) {
    const success = await this.webrtcService.switchCamera(deviceId);
    if (success) {
      this.selectedCameraId = deviceId;
    }
  }

  // Microphone change handler
  async onMicrophoneChange(deviceId: string) {
    const success = await this.webrtcService.switchMicrophone(deviceId);
    if (success) {
      this.selectedMicrophoneId = deviceId;
    }
  }

  // Toggle settings menu
  toggleSettings() {
    this.showSettings = !this.showSettings;
  }
}
