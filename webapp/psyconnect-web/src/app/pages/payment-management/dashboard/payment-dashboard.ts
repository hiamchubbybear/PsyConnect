import { Component } from '@angular/core';
import { ChartConfiguration, ChartOptions } from 'chart.js';

@Component({
  selector: 'app-dashboard',
  standalone : true,
  templateUrl: './payment-dashboard.html',
  styleUrls: ['./payment-dashboard.scss'],
})
export class DashboardComponent {

  chartOptions: ChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
        labels: { font: { size: 12 } },
      },
    },
    scales: {
      x: { grid: { display: false } },
      y: { grid: { color: '#eee' }, beginAtZero: true },
    },
  };

  // Income Report Chart
  incomeChartData: ChartConfiguration<'bar'>['data'] = {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul'],
    datasets: [
      {
        data: [1200, 1900, 3000, 2500, 3200, 4100, 3800],
        label: 'Income',
        backgroundColor: '#4caf50',
      },
    ],
  };

  // Spending Categories Chart
  spendingChartData: ChartConfiguration<'bar'>['data'] = {
    labels: ['Marketing', 'Operations', 'HR', 'Tech', 'Others'],
    datasets: [
      {
        data: [1500, 2000, 1200, 2800, 800],
        label: 'Spending',
        backgroundColor: '#f44336',
      },
    ],
  };

  // Transactions Data
  transactions = [
    {
      id: 'T001',
      customer: 'Nguyen Van A',
      method: 'Credit Card',
      amount: 250,
      status: 'Completed',
    },
    {
      id: 'T002',
      customer: 'Tran Thi B',
      method: 'Paypal',
      amount: 120,
      status: 'Pending',
    },
    {
      id: 'T003',
      customer: 'Pham Van C',
      method: 'Bank Transfer',
      amount: 560,
      status: 'Completed',
    },
    {
      id: 'T004',
      customer: 'Le Van D',
      method: 'Cash',
      amount: 75,
      status: 'Cancelled',
    },
  ];

  // Customers Data
  customers = [
    {
      name: 'Nguyen Van E',
      email: 'e.nguyen@example.com',
      amount: 1200,
      avatar: 'https://i.pravatar.cc/40?img=1',
    },
    {
      name: 'Tran Thi F',
      email: 'f.tran@example.com',
      amount: 850,
      avatar: 'https://i.pravatar.cc/40?img=2',
    },
    {
      name: 'Pham Van G',
      email: 'g.pham@example.com',
      amount: 640,
      avatar: 'https://i.pravatar.cc/40?img=3',
    },
  ];
}
