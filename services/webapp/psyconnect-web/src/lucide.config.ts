import { importProvidersFrom } from '@angular/core';
import { Check, ChevronLeft, ChevronRight, Eye, LucideAngularModule, MessageSquare, User, UserPlus, UserX, X } from 'lucide-angular';

export const lucideConfig = importProvidersFrom(
  LucideAngularModule.pick({ Eye, MessageSquare, UserPlus, X, Check,UserX , ChevronRight, ChevronLeft, User})
);
