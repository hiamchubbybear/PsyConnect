import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { delay } from 'rxjs/operators';

export interface MockPost {
  id: string;
  title: string;
  content: string;
  author_id: string;
  author_name: string;
  tags: string[];
  categories: string[];
  media: any[];
  visibility: 'public' | 'private' | 'followers';
  post_type: 'article' | 'question' | 'discussion' | 'resource';
  view_count: number;
  like_count: number;
  comment_count: number;
  share_count: number;
  status: 'published' | 'draft';
  created_at: string;
  updated_at: string;
}

export interface MockTransaction {
  id: string;
  type: 'payment' | 'refund' | 'subscription';
  amount: number;
  currency: string;
  status: 'completed' | 'pending' | 'failed';
  description: string;
  customer: string;
  date: string;
  method: string;
}

export interface MockInvoice {
  id: string;
  invoiceNumber: string;
  customer: string;
  amount: number;
  status: 'paid' | 'unpaid' | 'overdue';
  dueDate: string;
  items: { description: string; amount: number }[];
}

@Injectable({
  providedIn: 'root'
})
export class MockDataService {
  private posts: MockPost[] = [
    {
      id: '1',
      title: 'Understanding Anxiety: A Comprehensive Guide',
      content: 'Anxiety is a natural response to stress, but when it becomes overwhelming, it can significantly impact daily life. This guide explores the various types of anxiety disorders, their symptoms, and effective coping strategies. Learn about cognitive behavioral therapy, mindfulness techniques, and when to seek professional help.',
      author_id: 'user1',
      author_name: 'Dr. Sarah Johnson',
      tags: ['anxiety', 'mental-health', 'therapy'],
      categories: ['Mental Health'],
      media: [],
      visibility: 'public',
      post_type: 'article',
      view_count: 1245,
      like_count: 89,
      comment_count: 23,
      share_count: 45,
      status: 'published',
      created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: '2',
      title: 'Mindfulness Meditation for Beginners',
      content: 'Discover the transformative power of mindfulness meditation. This beginner-friendly guide covers basic techniques, breathing exercises, and tips for establishing a daily practice. Start your journey to inner peace and mental clarity today.',
      author_id: 'user2',
      author_name: 'Michael Chen',
      tags: ['mindfulness', 'meditation', 'wellness'],
      categories: ['Self Care'],
      media: [],
      visibility: 'public',
      post_type: 'article',
      view_count: 892,
      like_count: 67,
      comment_count: 15,
      share_count: 32,
      status: 'published',
      created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: '3',
      title: 'Building Healthy Relationships',
      content: 'Strong relationships are built on trust, communication, and mutual respect. Explore key principles for maintaining healthy connections with family, friends, and partners. Learn effective communication techniques and conflict resolution strategies.',
      author_id: 'user3',
      author_name: 'Emma Williams',
      tags: ['relationships', 'communication', 'wellness'],
      categories: ['Relationships'],
      media: [],
      visibility: 'public',
      post_type: 'discussion',
      view_count: 654,
      like_count: 45,
      comment_count: 28,
      share_count: 19,
      status: 'published',
      created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: '4',
      title: 'Coping with Depression: Practical Strategies',
      content: 'Depression affects millions worldwide. This article provides evidence-based strategies for managing depressive symptoms, including lifestyle changes, therapeutic approaches, and the importance of seeking professional support.',
      author_id: 'user1',
      author_name: 'Dr. Sarah Johnson',
      tags: ['depression', 'mental-health', 'coping'],
      categories: ['Mental Health'],
      media: [],
      visibility: 'public',
      post_type: 'article',
      view_count: 1567,
      like_count: 123,
      comment_count: 45,
      share_count: 67,
      status: 'published',
      created_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      id: '5',
      title: 'The Power of Positive Thinking',
      content: 'Cultivating a positive mindset can transform your life. Learn about the science behind positive psychology, practical exercises for reframing negative thoughts, and how optimism impacts mental and physical health.',
      author_id: 'user4',
      author_name: 'James Anderson',
      tags: ['positivity', 'mindset', 'wellness'],
      categories: ['Self Care'],
      media: [],
      visibility: 'public',
      post_type: 'article',
      view_count: 432,
      like_count: 34,
      comment_count: 12,
      share_count: 8,
      status: 'draft',
      created_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString()
    }
  ];

  private transactions: MockTransaction[] = [
    {
      id: 'txn1',
      type: 'payment',
      amount: 150.00,
      currency: 'USD',
      status: 'completed',
      description: 'Therapy Session - Dr. Sarah Johnson',
      customer: 'John Doe',
      date: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString(),
      method: 'Credit Card'
    },
    {
      id: 'txn2',
      type: 'subscription',
      amount: 29.99,
      currency: 'USD',
      status: 'completed',
      description: 'Monthly Premium Subscription',
      customer: 'Jane Smith',
      date: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
      method: 'PayPal'
    },
    {
      id: 'txn3',
      type: 'payment',
      amount: 200.00,
      currency: 'USD',
      status: 'pending',
      description: 'Group Therapy Session',
      customer: 'Mike Johnson',
      date: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
      method: 'Bank Transfer'
    },
    {
      id: 'txn4',
      type: 'refund',
      amount: -75.00,
      currency: 'USD',
      status: 'completed',
      description: 'Cancelled Session Refund',
      customer: 'Sarah Williams',
      date: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
      method: 'Credit Card'
    }
  ];

  private invoices: MockInvoice[] = [
    {
      id: 'inv1',
      invoiceNumber: 'INV-2024-001',
      customer: 'John Doe',
      amount: 450.00,
      status: 'paid',
      dueDate: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000).toISOString(),
      items: [
        { description: 'Therapy Session (3x)', amount: 450.00 }
      ]
    },
    {
      id: 'inv2',
      invoiceNumber: 'INV-2024-002',
      customer: 'Jane Smith',
      amount: 300.00,
      status: 'unpaid',
      dueDate: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString(),
      items: [
        { description: 'Therapy Session (2x)', amount: 300.00 }
      ]
    },
    {
      id: 'inv3',
      invoiceNumber: 'INV-2024-003',
      customer: 'Mike Johnson',
      amount: 600.00,
      status: 'overdue',
      dueDate: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000).toISOString(),
      items: [
        { description: 'Therapy Session (4x)', amount: 600.00 }
      ]
    }
  ];

  getPosts(): Observable<MockPost[]> {
    return of(this.posts).pipe(delay(300));
  }

  getPostById(id: string): Observable<MockPost | undefined> {
    return of(this.posts.find(p => p.id === id)).pipe(delay(200));
  }

  getTransactions(): Observable<MockTransaction[]> {
    return of(this.transactions).pipe(delay(300));
  }

  getInvoices(): Observable<MockInvoice[]> {
    return of(this.invoices).pipe(delay(300));
  }

  
  getPaymentStats() {
    const totalRevenue = this.transactions
      .filter(t => t.status === 'completed' && t.amount > 0)
      .reduce((sum, t) => sum + t.amount, 0);

    const monthlyRevenue = this.transactions
      .filter(t => {
        const txnDate = new Date(t.date);
        const now = new Date();
        return txnDate.getMonth() === now.getMonth() &&
               txnDate.getFullYear() === now.getFullYear() &&
               t.status === 'completed' && t.amount > 0;
      })
      .reduce((sum, t) => sum + t.amount, 0);

    return of({
      totalRevenue,
      monthlyRevenue,
      totalTransactions: this.transactions.length,
      pendingPayments: this.transactions.filter(t => t.status === 'pending').length
    }).pipe(delay(200));
  }
}
