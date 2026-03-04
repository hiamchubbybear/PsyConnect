import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { ConsultationSession } from '../../../models/consultation.model';

@Component({
  selector: 'app-session-card',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="session-card" [class.paid]="session.payment_status === 'PAID'">
      <div class="session-card__header">
        <div class="session-card__icon-wrap">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M8 2v4"></path>
            <path d="M16 2v4"></path>
            <rect width="18" height="18" x="3" y="4" rx="2"></rect>
            <path d="M3 10h18"></path>
          </svg>
        </div>
        <div class="session-card__title-wrap">
          <h4>
            {{
              session.mode === 'online'
                ? 'Tư vấn Trực tuyến'
                : 'Tư vấn Trực tiếp'
            }}
          </h4>
          <span
            class="payment-badge"
            [class.paid]="session.payment_status === 'PAID'"
          >
            {{
              session.payment_status === 'PAID'
                ? 'Đã thanh toán'
                : 'Chờ thanh toán'
            }}
          </span>
        </div>
      </div>

      <div class="session-card__body">
        <div class="detail-item">
          <div class="detail-item__icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
          </div>
          <div class="detail-item__content">
            <span class="label">Ngày hẹn</span>
            <span class="value">{{
              session.scheduled_date || session.start_time | date: 'fullDate'
            }}</span>
          </div>
        </div>

        <div class="detail-item">
          <div class="detail-item__icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </div>
          <div class="detail-item__content">
            <span class="label">Thời gian</span>
            <span class="value"
              >{{ session.start_time | date: 'shortTime' }} -
              {{ session.end_time | date: 'shortTime' }}</span
            >
          </div>
        </div>

        <div
          class="detail-item"
          *ngIf="
            session.location_info?.address_line && session.mode !== 'online'
          "
        >
          <div class="detail-item__icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
              <circle cx="12" cy="10" r="3"></circle>
            </svg>
          </div>
          <div class="detail-item__content">
            <span class="label">Địa điểm</span>
            <span class="value address">{{
              session.location_info?.address_line
            }}</span>
          </div>
        </div>
      </div>

      <div class="session-card__footer">
        <button class="action-btn action-btn--primary">
          {{
            session.payment_status === 'PAID'
              ? 'Xem chi tiết'
              : 'Thanh toán ngay'
          }}
        </button>
      </div>
    </div>
  `,
  styles: [
    `
      .session-card {
        background: #ffffff;
        border: 1px solid var(--gray-200);
        border-radius: 16px;
        width: 300px;
        margin-top: 10px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
        overflow: hidden;
        transition:
          transform 0.2s,
          box-shadow 0.2s;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
        }

        &.paid {
          border-left: 4px solid var(--secondary-color);
        }
      }

      .session-card__header {
        padding: 16px;
        display: flex;
        align-items: center;
        gap: 12px;
        background: var(--gray-50);
        border-bottom: 1px solid var(--gray-100);
      }

      .session-card__icon-wrap {
        width: 36px;
        height: 36px;
        background: var(--primary-color);
        color: white;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .session-card__title-wrap {
        flex: 1;
        h4 {
          margin: 0 0 2px;
          font-size: 0.95rem;
          font-weight: 700;
          color: var(--gray-900);
        }
      }

      .payment-badge {
        font-size: 0.7rem;
        font-weight: 600;
        padding: 2px 8px;
        border-radius: 12px;
        background: #fee2e2;
        color: #ef4444;

        &.paid {
          background: #dcfce7;
          color: #10b981;
        }
      }

      .session-card__body {
        padding: 16px;
        display: flex;
        flex-direction: column;
        gap: 12px;
      }

      .detail-item {
        display: flex;
        gap: 10px;
        align-items: flex-start;
      }

      .detail-item__icon {
        color: var(--gray-400);
        margin-top: 2px;
      }

      .detail-item__content {
        display: flex;
        flex-direction: column;
        .label {
          font-size: 0.7rem;
          color: var(--gray-400);
          text-transform: uppercase;
          letter-spacing: 0.05em;
        }
        .value {
          font-size: 0.85rem;
          font-weight: 500;
          color: var(--gray-700);
          &.address {
            color: var(--primary-color);
            text-decoration: underline;
            cursor: pointer;
          }
        }
      }

      .session-card__footer {
        padding: 12px 16px 16px;
      }

      .action-btn {
        width: 100%;
        padding: 10px;
        border-radius: 10px;
        font-size: 0.85rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
        border: none;

        &--primary {
          background: var(--primary-color);
          color: white;
          &:hover {
            background: var(--primary-dark);
          }
        }
      }
    `,
  ],
})
export class SessionCardComponent {
  @Input() session!: ConsultationSession;
}
