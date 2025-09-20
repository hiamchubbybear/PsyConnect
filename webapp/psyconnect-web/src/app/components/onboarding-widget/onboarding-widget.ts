// onboarding-modal.interface.ts
export interface OnboardingStep {
  id: string;
  icon?: string; // Icon class or URL
  title: string;
  description: string;
  primaryButton: OnboardingButton;
  secondaryButton?: OnboardingButton;
  backgroundGradient?: string;
  image?: string; // Optional image URL
}

export interface OnboardingButton {
  text: string;
  action?: () => void;
  disabled?: boolean;
  loading?: boolean;
  icon?: string;
}

export interface OnboardingConfig {
  steps: OnboardingStep[];
  showProgressDots?: boolean;
  showCloseButton?: boolean;
  backdropClosable?: boolean;
  theme?: 'light' | 'dark';
  onComplete?: () => void;
  onClose?: () => void;
}

import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-onboarding-modal',
  imports: [CommonModule],
  templateUrl: './onboarding-widget.html',
  styleUrls: ['./onboarding-widget.scss'],
})
export class OnboardingModalComponent implements OnInit {
  @Input() config!: OnboardingConfig;
  @Input() isVisible: boolean = false;
  @Output() close = new EventEmitter<void>();
  @Output() complete = new EventEmitter<void>();
  @Output() stepChange = new EventEmitter<{
    step: OnboardingStep;
    index: number;
  }>();

  currentStepIndex = 0;
  isAnimating = false;

  ngOnInit() {
    if (!this.config || !this.config.steps.length) {
      console.error('OnboardingModal: config with steps is required');
      return;
    }
  }

  get currentStep(): OnboardingStep {
    return this.config.steps[this.currentStepIndex];
  }

  get isFirstStep(): boolean {
    return this.currentStepIndex === 0;
  }

  get isLastStep(): boolean {
    return this.currentStepIndex === this.config.steps.length - 1;
  }

  get progressPercentage(): number {
    return ((this.currentStepIndex + 1) / this.config.steps.length) * 100;
  }

  onBackdropClick(): void {
    if (this.config.backdropClosable !== false) {
      this.onClose();
    }
  }

  onClose(): void {
    this.config.onClose?.();
    this.close.emit();
  }

  onPrimaryAction(): void {
    if (this.currentStep.primaryButton.disabled || this.isAnimating) {
      return;
    }

    if (this.currentStep.primaryButton.action) {
      this.currentStep.primaryButton.action();
    }

    if (this.isLastStep) {
      this.onComplete();
    } else {
      this.nextStep();
    }
  }

  onSecondaryAction(): void {
    if (
      !this.currentStep.secondaryButton ||
      this.currentStep.secondaryButton.disabled ||
      this.isAnimating
    ) {
      return;
    }

    if (this.currentStep.secondaryButton.action) {
      this.currentStep.secondaryButton.action();
    } else {
      this.prevStep();
    }
  }

  nextStep(): void {
    if (this.isLastStep || this.isAnimating) return;

    this.isAnimating = true;
    setTimeout(() => {
      this.currentStepIndex++;
      this.stepChange.emit({
        step: this.currentStep,
        index: this.currentStepIndex,
      });
      this.isAnimating = false;
    }, 150);
  }

  prevStep(): void {
    if (this.isFirstStep || this.isAnimating) return;

    this.isAnimating = true;
    setTimeout(() => {
      this.currentStepIndex--;
      this.stepChange.emit({
        step: this.currentStep,
        index: this.currentStepIndex,
      });
      this.isAnimating = false;
    }, 150);
  }

  goToStep(index: number): void {
    if (
      index < 0 ||
      index >= this.config.steps.length ||
      index === this.currentStepIndex ||
      this.isAnimating
    ) {
      return;
    }

    this.isAnimating = true;
    setTimeout(() => {
      this.currentStepIndex = index;
      this.stepChange.emit({
        step: this.currentStep,
        index: this.currentStepIndex,
      });
      this.isAnimating = false;
    }, 150);
  }

  onComplete(): void {
    this.config.onComplete?.();
    this.complete.emit();
  }

  // Utility methods for template
  getBackgroundStyle(): any {
    return {
      background:
        this.currentStep.backgroundGradient ||
        'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    };
  }

  getThemeClass(): string {
    return this.config.theme === 'dark' ? 'theme-dark' : 'theme-light';
  }
}
