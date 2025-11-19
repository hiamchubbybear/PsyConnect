import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ThemeService } from '../../services/theme/theme-service';
import { AudienceSectionComponent } from './about-adience';
import { ContactSectionComponent } from './about-contact';
import { DevelopmentSectionComponent } from './about-development';
import { FaqSectionComponent } from './about-faq-section';
import { HeroSectionComponent } from './about-hero-section';
import { PartnersSectionComponent } from "./about-partners-section";
import { PrivacySectionComponent } from './about-privacy-section';
import { ServicesSectionComponent } from "./about-service-section";
import { StorySectionComponent } from "./about-story-section";
import { TeamSectionComponent } from "./about-team-section";
import { AboutUsDataService } from './about-us-data-service';
import { VerificationSectionComponent } from './about-verification';


@Component({
  selector: 'app-about-us',
  standalone: true,
  imports: [
    TranslateModule,
    CommonModule,
    RouterModule,
    HeroSectionComponent,
    VerificationSectionComponent,
    AudienceSectionComponent,
    HeroSectionComponent,
    PrivacySectionComponent,
    FaqSectionComponent,
    ContactSectionComponent,
    DevelopmentSectionComponent,
    CommonModule,
    PartnersSectionComponent,
    ServicesSectionComponent,
    TeamSectionComponent,
    StorySectionComponent
],
  template: `
    <main class="about-us">
      <app-hero-section
        [brand]="brand"
        [statistics]="statistics"
        [mission]="mission"
        [coreValues]="coreValues"
        (bookNow)="onBookNow()"
        (listAsPartner)="onListAsPartner()">
      </app-hero-section>

      <app-story-section
        [storyTitle]="storyTitle"
        [story]="story"
        [mission]="mission"
        [vision]="vision">
      </app-story-section>

      <app-verification-section
        [verificationSteps]="verificationSteps"
        [qualityStandards]="qualityStandards">
      </app-verification-section>

      <app-partners-section
        [partners]="partners"
        (verifyCertificate)="verifyCertificate($event)">
      </app-partners-section>

      <app-services-section
        [services]="services"
        [methods]="methods">
      </app-services-section>

      <app-audience-section
        [targetAudience]="targetAudience">
      </app-audience-section>

      <app-team-section
        [team]="team">
      </app-team-section>

      <app-privacy-section
        [privacyFeatures]="privacyFeatures"
        [securityMeasures]="securityMeasures">
      </app-privacy-section>

      <app-faq-section
        [faq]="faq">
      </app-faq-section>

      <app-contact-section
        [contact]="contact"
        (emergencyHelp)="onEmergencyHelp()">
      </app-contact-section>

      <app-development-section
        [devs]="devs"
        [release]="release"
        (viewRepo)="onViewRepo()">
      </app-development-section>
    </main>
  `,
  styleUrls: ['./about-us.scss'],
})
export class AboutUsComponent implements OnInit {
   brand: any;
  statistics: any;
  partners: any;
  services: any;
  team: any;
  contact: any;
  devs: any;
  release: any;

  constructor(
    private translate: TranslateService,
    private themeService: ThemeService,
    private dataService: AboutUsDataService
  ) {}

  ngOnInit() {
    this.themeService.setTheme('light');
    this.loadTranslations();

    this.brand = this.dataService.getBrandInfo();
    this.statistics = this.dataService.getStatistics();
    this.partners = this.dataService.getPartners();
    this.services = this.dataService.getServices();
    this.team = this.dataService.getTeam();
    this.contact = this.dataService.getContact();
    this.devs = this.dataService.getDevs();
    this.release = this.dataService.getRelease();
  }
  mission = '';
  vision = '';
  coreValues: string[] = [];
  storyTitle = '';
  story = '';
  verificationSteps: string[] = [];
  qualityStandards: string[] = [];
  privacyFeatures: string[] = [];
  securityMeasures: string[] = [];
  methods: string[] = [];
  targetAudience: string[] = [];
  faq: { question: string; answer: string }[] = [];

  loadTranslations() {
    this.translate
      .get([
        'mission_statement',
        'vision_statement',
        'core_values',
        'our_story_title',
        'our_story_content',
        'verification_steps',
        'quality_standards',
        'privacy_features',
        'security_measures',
        'therapeutic_methods',
        'target_audience',
        'faq_items',
      ])
      .subscribe((translations) => {
        this.mission = translations['mission_statement'];
        this.vision = translations['vision_statement'];
        this.coreValues = translations['core_values'];
        this.storyTitle = translations['our_story_title'];
        this.story = translations['our_story_content'];
        this.verificationSteps = translations['verification_steps'];
        this.qualityStandards = translations['quality_standards'];
        this.privacyFeatures = translations['privacy_features'];
        this.securityMeasures = translations['security_measures'];
        this.methods = translations['therapeutic_methods'];
        this.targetAudience = translations['target_audience'];
        this.faq = translations['faq_items'];
      });
  }

  onBookNow() {
    window.location.href = '/search';
  }

  onListAsPartner() {
    window.location.href = '/partner/signup';
  }

  onViewRepo() {
    window.open('https://github.com/hiamchubbybear/psyconnect', '_blank');
  }

  onEmergencyHelp() {
    window.location.href = `tel:${this.contact.emergency}`;
  }

  verifyCertificate(cert: any) {
    if (cert.verificationUrl) {
      window.open(cert.verificationUrl, '_blank');
    }
  }
}
