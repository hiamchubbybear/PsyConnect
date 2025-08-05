import { animate, style, transition, trigger } from '@angular/animations';

export const fadeRouteAnimation = trigger('fadeAnimation', [
  transition('* <=> *', [
    style({ opacity: 0 }),
    animate('0.2s ease', style({ opacity: 1 })),
  ]),
]);
