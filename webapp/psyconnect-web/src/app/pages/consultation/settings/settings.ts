import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
    OnboardingConfig,
    OnboardingModalComponent,
    OnboardingStep,
} from '../../../components/onboarding-widget/onboarding-widget';

@Component({
  selector: 'consultation-settings',
  standalone: true,
  imports: [CommonModule, OnboardingModalComponent],
  templateUrl: './settings.html',
  styleUrls: ['./settings.scss'],
})
export class Settings {
  showOnboarding = true;

  onboardingConfig: OnboardingConfig = {
    steps: [
      {
        id: 'welcome',
        title: 'Welcome',
        description: 'Let’s set up your preferences before you start.',
        backgroundGradient: 'linear-gradient(135deg, #000000 0%, #1a1a1a 100%)',
        primaryButton: {
          text: 'Continue',
          action: () => console.log('Welcome step completed'),
        },
        secondaryButton: {
          text: 'Skip',
        },
      },
      {
        id: 'features',
        title: 'What You Can Do',
        description:
          'Quickly explore the core features designed to support your journey.',
        backgroundGradient: 'linear-gradient(135deg, #1a1a1a 0%, #333333 100%)',
        primaryButton: {
          text: 'Next',
          icon: 'fas fa-arrow-right',
        },
        secondaryButton: {
          text: 'Back',
          icon: 'fas fa-arrow-left',
        },
      },
      {
        id: 'complete',
        title: 'You’re All Set',
        description:
          'Everything is ready. Start using the app with confidence.',
        backgroundGradient: 'linear-gradient(135deg, #0d0d0d 0%, #1f1f1f 100%)',
        primaryButton: {
          text: 'Start Now',
          action: () => console.log('Onboarding finished'),
        },
        secondaryButton: {
          text: 'Learn More',
          icon: 'fas fa-info-circle',
          action: () => window.open('https://example.com/help', '_blank'),
        },
      },
    ],

    showProgressDots: true,
    showCloseButton: true,
    backdropClosable: true,
    theme: 'light',

    onComplete: () => console.log('Onboarding completed!'),
    onClose: () => console.log('Onboarding closed'),
  };

  onOnboardingComplete(): void {
    this.showOnboarding = false;
    alert('Welcome aboard! ');
  }

  onStepChange(event: { step: OnboardingStep; index: number }): void {
    console.log('Step changed:', event.step.title, 'Index:', event.index);
  }
}
