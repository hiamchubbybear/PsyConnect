import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';

interface OnboardingStep {
  id: number;
  title: string;
  description: string;
  image: string;
  mascot: 'both' | 'psy' | 'connect';
}

@Component({
  selector: 'app-onboarding',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './onboarding.component.html',
  styleUrls: ['./onboarding.component.scss'],
})
export class OnboardingComponent {
  @Output() completed = new EventEmitter<void>();
  @Output() skipped = new EventEmitter<void>();

  currentStep = 0;

  steps: OnboardingStep[] = [
    {
      id: 1,
      title: 'Chào mừng đến với PsyConnect!',
      description:
        'Nền tảng kết nối bạn với các chuyên gia tâm lý và cộng đồng hỗ trợ',
      image: 'assets/mascot/expressions/happy.png',
      mascot: 'both',
    },
    {
      id: 2,
      title: 'Trò chuyện với emotes đáng yêu',
      description:
        'Thể hiện cảm xúc của bạn với hơn 20 emotes dễ thương được thiết kế riêng',
      image: 'assets/emotes-svg/love.svg',
      mascot: 'psy',
    },
    {
      id: 3,
      title: 'Kết nối với chuyên gia',
      description:
        'Tìm và đặt lịch với các nhà tâm lý chuyên nghiệp phù hợp với bạn',
      image: 'assets/mascot/connect-blue.png',
      mascot: 'connect',
    },
    {
      id: 4,
      title: 'Tham gia cộng đồng',
      description: 'Chia sẻ, học hỏi và nhận hỗ trợ từ cộng đồng thân thiện',
      image: 'assets/illustrations/illustration_success.png',
      mascot: 'both',
    },
    {
      id: 5,
      title: 'Sẵn sàng bắt đầu!',
      description:
        'Hãy cùng nhau tạo nên một hành trình chăm sóc sức khỏe tinh thần tích cực',
      image: 'assets/mascot/expressions/happy.png',
      mascot: 'both',
    },
  ];

  get currentStepData(): OnboardingStep {
    return this.steps[this.currentStep];
  }

  get isFirstStep(): boolean {
    return this.currentStep === 0;
  }

  get isLastStep(): boolean {
    return this.currentStep === this.steps.length - 1;
  }

  get progress(): number {
    return ((this.currentStep + 1) / this.steps.length) * 100;
  }

  nextStep() {
    if (this.isLastStep) {
      this.complete();
    } else {
      this.currentStep++;
    }
  }

  previousStep() {
    if (!this.isFirstStep) {
      this.currentStep--;
    }
  }

  skip() {
    this.markAsCompleted();
    this.skipped.emit();
  }

  complete() {
    this.markAsCompleted();
    this.completed.emit();
  }

  private markAsCompleted() {
    localStorage.setItem('onboardingCompleted', 'true');
  }

  static hasCompletedOnboarding(): boolean {
    return localStorage.getItem('onboardingCompleted') === 'true';
  }
}
