import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    EventEmitter,
    HostListener,
    Input,
    OnDestroy,
    OnInit,
    Output,
    forwardRef,
} from '@angular/core';
import {
    ControlValueAccessor,
    FormsModule,
    NG_VALUE_ACCESSOR,
} from '@angular/forms';

export interface DropdownOption {
  value: any;
  label: string;
  disabled?: boolean;
  icon?: string;
}

@Component({
  selector: 'app-dropdown',
  standalone: true,
  imports: [CommonModule, FormsModule],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => DropdownComponent),
      multi: true,
    },
  ],
  template: `
    <div class="dropdown" [class.dropdown--disabled]="disabled">
      <!-- Trigger Button -->
      <button
        type="button"
        class="dropdown__trigger"
        [class.dropdown__trigger--open]="isOpen"
        [class.dropdown__trigger--error]="hasError"
        [disabled]="disabled"
        (click)="toggle()"
        (keydown.arrowdown)="onArrowDown($event)"
        (keydown.arrowup)="onArrowUp($event)"
        (keydown.enter)="onEnter($event)"
        (keydown.space)="onSpace($event)"
        [attr.aria-expanded]="isOpen"
        [attr.aria-haspopup]="'listbox'"
        [attr.aria-label]="ariaLabel"
      >
        <span class="dropdown__value">
          <span *ngIf="selectedOption?.icon" class="dropdown__icon">{{
            selectedOption?.icon
          }}</span>
          {{ selectedOption?.label || placeholder }}
        </span>
        <span class="dropdown__arrow" [class.dropdown__arrow--open]="isOpen">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path
              d="M4.427 6.573L8 10.146l3.573-3.573a.5.5 0 11.707.708l-3.927 3.926a.5.5 0 01-.707 0L3.72 7.28a.5.5 0 11.707-.707z"
            />
          </svg>
        </span>
      </button>

      <!-- Error Message -->
      <div *ngIf="hasError && errorMessage" class="dropdown__error">
        {{ errorMessage }}
      </div>

      <!-- Dropdown Menu -->
      <div
        class="dropdown__menu"
        [class.dropdown__menu--open]="isOpen"
        [style.min-width.px]="minWidth"
        role="listbox"
        [attr.aria-label]="ariaLabel"
      >
        <div
          *ngFor="let option of filteredOptions; let i = index"
          class="dropdown__option"
          [class.dropdown__option--selected]="option.value === selectedValue"
          [class.dropdown__option--focused]="i === focusedIndex"
          [class.dropdown__option--disabled]="option.disabled"
          role="option"
          [attr.aria-selected]="option.value === selectedValue"
          (click)="selectOption(option)"
          (mouseenter)="focusedIndex = i"
        >
          <span *ngIf="option.icon" class="dropdown__option-icon">{{
            option.icon
          }}</span>
          {{ option.label }}
        </div>

        <!-- Empty State -->
        <div *ngIf="filteredOptions.length === 0" class="dropdown__empty">
          {{ emptyText }}
        </div>
      </div>

      <!-- Loading State -->
      <div *ngIf="loading" class="dropdown__loading">
        <div class="dropdown__spinner"></div>
      </div>
    </div>
  `,
  styleUrls: ['./dropdown.scss'],
})
export class DropdownComponent
  implements OnInit, OnDestroy, ControlValueAccessor
{
  @Input() options: DropdownOption[] = [];
  @Input() placeholder: string = 'Select an option';
  @Input() disabled: boolean = false;
  @Input() loading: boolean = false;
  @Input() hasError: boolean = false;
  @Input() errorMessage: string = '';
  @Input() emptyText: string = 'No options available';
  @Input() minWidth: number = 200;
  @Input() ariaLabel: string = 'Dropdown menu';

  @Output() selectionChange = new EventEmitter<DropdownOption | null>();
  @Output() opened = new EventEmitter<void>();
  @Output() closed = new EventEmitter<void>();
  private _options: DropdownOption[] = [];
  isOpen = false;
  selectedValue: any = null;
  selectedOption: DropdownOption | null = null;
  focusedIndex = -1;
  filteredOptions: DropdownOption[] = [];

  @Input() set setOptions(val: DropdownOption[]) {
    this._options = val;
    this.filteredOptions = this._options.filter((option) => !option.disabled);
  }
  get getOptions() {
    return this._options;
  }
  private onChange = (value: any) => {};
  private onTouched = () => {};

  constructor(private elementRef: ElementRef) {}

  ngOnInit() {
    this.filteredOptions = this.options.filter((option) => !option.disabled);
  }

  ngOnDestroy() {
    this.close();
  }

  
  writeValue(value: any): void {
    this.selectedValue = value;
    this.selectedOption =
      this.options.find((option) => option.value === value) || null;
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  
  toggle(): void {
    if (this.disabled || this.loading) return;
    this.isOpen ? this.close() : this.open();
  }

  open(): void {
    if (this.disabled || this.loading) return;
    this.isOpen = true;
    this.focusedIndex = this.selectedOption
      ? this.filteredOptions.findIndex(
          (option) => option.value === this.selectedValue
        )
      : 0;
    this.opened.emit();
  }

  close(): void {
    this.isOpen = false;
    this.focusedIndex = -1;
    this.onTouched();
    this.closed.emit();
  }

  selectOption(option: DropdownOption): void {
    if (option.disabled) return;

    this.selectedValue = option.value;
    this.selectedOption = option;
    this.onChange(option.value);
    this.selectionChange.emit(option);
    this.close();
  }

  onArrowDown(event: Event): void {
    const keyboardEvent = event as KeyboardEvent;
    keyboardEvent.preventDefault();
    if (!this.isOpen) {
      this.open();
      return;
    }
    this.focusedIndex = Math.min(
      this.focusedIndex + 1,
      this.filteredOptions.length - 1
    );
  }

  onArrowUp(event: Event): void {
    const keyboardEvent = event as KeyboardEvent;
    keyboardEvent.preventDefault();
    if (!this.isOpen) {
      this.open();
      return;
    }
    this.focusedIndex = Math.max(this.focusedIndex - 1, 0);
  }

  onEnter(event: Event): void {
    const keyboardEvent = event as KeyboardEvent;
    keyboardEvent.preventDefault();
    if (this.isOpen && this.focusedIndex >= 0) {
      this.selectOption(this.filteredOptions[this.focusedIndex]);
    } else {
      this.toggle();
    }
  }

  onSpace(event: Event): void {
    const keyboardEvent = event as KeyboardEvent;
    keyboardEvent.preventDefault();
    this.toggle();
  }

  
  @HostListener('document:click', ['$event'])
  onDocumentClick(event: Event): void {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.close();
    }
  }

  @HostListener('document:keydown.escape')
  onEscapeKey(): void {
    this.close();
  }
}
