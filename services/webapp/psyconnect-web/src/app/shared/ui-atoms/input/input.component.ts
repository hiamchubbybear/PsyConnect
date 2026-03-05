import { CommonModule } from '@angular/common';
import { Component, Input, OnInit, forwardRef } from '@angular/core';
import {
  ControlValueAccessor,
  FormControl,
  NG_VALUE_ACCESSOR,
  ReactiveFormsModule,
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-input',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './input.component.html',
  styleUrls: ['./input.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => InputComponent),
      multi: true,
    },
  ],
})
export class InputComponent implements ControlValueAccessor, OnInit {
  @Input() label?: string;
  @Input() placeholder = '';
  @Input() type = 'text';
  @Input() icon?: string; 
  @Input() errorMessages: Record<string, string> = {};
  @Input() control?: FormControl;

  
  value: any = '';
  isDisabled = false;
  showPassword = false;

  protected readonly objectKeys = Object.keys;

  
  onChange = (value: any) => {};
  onTouched = () => {};

  ngOnInit() {
    
    
  }

  
  writeValue(value: any): void {
    this.value = value;
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.isDisabled = isDisabled;
  }

  
  onInput(event: Event) {
    const target = event.target as HTMLInputElement;
    this.value = target.value;
    this.onChange(this.value);
  }

  togglePasswordVisibility() {
    this.showPassword = !this.showPassword;
  }

  get inputType(): string {
    if (this.type === 'password') {
      return this.showPassword ? 'text' : 'password';
    }
    return this.type;
  }

  get isInvalid(): boolean {
    return !!(this.control && this.control.invalid && this.control.touched);
  }
}
