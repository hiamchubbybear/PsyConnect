-- =========================================
-- ALL-IN-ONE SQL: RESET AND CREATE HR + PAYROLL + INTEGRATED
-- Chạy 1 lần để xóa và tạo lại toàn bộ DB
-- =========================================

-- =========================================
-- DATABASE HR
-- =========================================
DROP DATABASE IF EXISTS hr;
CREATE DATABASE hr;
USE hr;

-- Benefits Plans (FK employee)
CREATE TABLE benefits_plan (
  benefits_plan_id INT AUTO_INCREMENT PRIMARY KEY,
  plan_name VARCHAR(100) NOT NULL,
  default_cost DECIMAL(12,2) DEFAULT 0.00,
  notes TEXT
);

-- Locations
CREATE TABLE location (
  location_id INT AUTO_INCREMENT PRIMARY KEY,
  location_name VARCHAR(255) NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Branches
CREATE TABLE branch (
  branch_id INT AUTO_INCREMENT PRIMARY KEY,
  branch_name VARCHAR(255) NOT NULL,
  location_id INT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (location_id) REFERENCES location(location_id) ON DELETE SET NULL
);

-- Departments
CREATE TABLE department (
  department_id INT AUTO_INCREMENT PRIMARY KEY,
  department_name VARCHAR(255) NOT NULL,
  number_of_employees INT DEFAULT 0,
  phone VARCHAR(50),
  address VARCHAR(255),
  email VARCHAR(255),
  notes TEXT,
  branch_id INT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (branch_id) REFERENCES branch(branch_id) ON DELETE SET NULL
);

-- Positions
CREATE TABLE job_position (
  position_id INT AUTO_INCREMENT PRIMARY KEY,
  position_name VARCHAR(255) NOT NULL
);

-- Employees
CREATE TABLE employee (
  employee_id CHAR(36) PRIMARY KEY,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(50),
  address VARCHAR(255),
  gender ENUM('Male','Female','Other') DEFAULT 'Other',
  ethnicity VARCHAR(100),
  birthday DATE,
  hire_date DATE,
  employment_type ENUM('Full-time','Part-time') DEFAULT 'Full-time',
  is_shareholder TINYINT(1) DEFAULT 0,
  salary DECIMAL(15,2) DEFAULT 0.00,
  vacation_days_used INT DEFAULT 0,
  vacation_days_total INT DEFAULT 0,
  benefits_plan_id INT,
  benefits_cost DECIMAL(12,2) DEFAULT 0.00,
  birthday_reminder_month TINYINT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  department_id INT,
  CONSTRAINT fk_emp_dept FOREIGN KEY (department_id) REFERENCES department(department_id) ON DELETE SET NULL,
  CONSTRAINT fk_emp_benefit FOREIGN KEY (benefits_plan_id) REFERENCES benefits_plan(benefits_plan_id)
);

-- Employee <-> Position
CREATE TABLE employee_position (
  id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  position_id INT NOT NULL,
  start_date DATE,
  end_date DATE,
  is_primary TINYINT(1) DEFAULT 0,
  FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE,
  FOREIGN KEY (position_id) REFERENCES job_position(position_id) ON DELETE CASCADE
);

-- Education
CREATE TABLE education (
  education_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  education_name VARCHAR(255),
  school VARCHAR(255),
  start_date DATE,
  end_date DATE,
  FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

-- Holidays
CREATE TABLE holidays (
  dayoff_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36),
  dayoff_date DATE,
  number_of_days DECIMAL(5,2) DEFAULT 1,
  reason VARCHAR(255),
  FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_employee_email ON employee(email);
CREATE INDEX idx_employee_department ON employee(department_id);
CREATE INDEX idx_employee_hiredate ON employee(hire_date);
CREATE INDEX idx_employee_birthday_month ON employee(birthday);


-- =========================================
-- DATABASE PAYROLL
-- =========================================
DROP DATABASE IF EXISTS payroll;
CREATE DATABASE payroll;
USE payroll;

-- Basic Salary
CREATE TABLE basic_salary (
  basicsalary_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  coefficient DECIMAL(8,4) DEFAULT 1.0,
  base_amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'USD',
  effective_date DATE NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Taxes
CREATE TABLE tax (
  tax_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  tax_type VARCHAR(100),
  tax_amount DECIMAL(12,2) NOT NULL,
  tax_date DATE,
  notes TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Penalties
CREATE TABLE penalty (
  fines_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  fines_type VARCHAR(255),
  fine_date DATE,
  fines_amount DECIMAL(12,2),
  reason TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Bonus
CREATE TABLE bonus (
  bonus_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  bonus_type VARCHAR(100),
  amount DECIMAL(12,2) NOT NULL,
  reason TEXT,
  bonus_date DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Allowances
CREATE TABLE allowance (
  allowance_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  allowance_type VARCHAR(100),
  allowance_amount DECIMAL(12,2) DEFAULT 0.00,
  effective_date DATE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Overtime
CREATE TABLE overtime (
  overtime_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  hours DECIMAL(6,2) DEFAULT 0,
  rate DECIMAL(10,2) DEFAULT 0,
  overtime_date DATE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Salary Payments
CREATE TABLE salary_payment (
  salary_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  payment_date DATE NOT NULL,
  payment_method VARCHAR(50),
  salary_amount DECIMAL(15,2) NOT NULL,
  basicsalary_id INT,
  overtime_id INT,
  dayoff_id INT,
  bonus_id INT,
  tax_id INT,
  allowance_id INT,
  fines_id INT,
  notes TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (basicsalary_id) REFERENCES basic_salary(basicsalary_id) ON DELETE SET NULL,
  FOREIGN KEY (overtime_id) REFERENCES overtime(overtime_id) ON DELETE SET NULL,
  FOREIGN KEY (bonus_id) REFERENCES bonus(bonus_id) ON DELETE SET NULL,
  FOREIGN KEY (tax_id) REFERENCES tax(tax_id) ON DELETE SET NULL,
  FOREIGN KEY (allowance_id) REFERENCES allowance(allowance_id) ON DELETE SET NULL,
  FOREIGN KEY (fines_id) REFERENCES penalty(fines_id) ON DELETE SET NULL
);

-- Vacation Records
CREATE TABLE vacation_record (
  vacation_id INT AUTO_INCREMENT PRIMARY KEY,
  employee_id CHAR(36) NOT NULL,
  year INT,
  vacation_total INT DEFAULT 0,
  vacation_used INT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- DATABASE INTEGRATED
-- =========================================
DROP DATABASE IF EXISTS integrated;
CREATE DATABASE integrated;
USE integrated;

-- Users
CREATE TABLE app_user (
  user_id CHAR(36) PRIMARY KEY,
  username VARCHAR(150) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  display_name VARCHAR(255),
  is_active TINYINT(1) DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_login DATETIME,
  employee_id CHAR(36) NULL
);

-- Roles
CREATE TABLE role (
  role_id INT AUTO_INCREMENT PRIMARY KEY,
  role_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE user_role (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  role_id INT NOT NULL,
  FOREIGN KEY (role_id) REFERENCES role(role_id) ON DELETE CASCADE
);

-- Refresh Tokens
CREATE TABLE refresh_token (
  token_id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  token_value VARCHAR(512) NOT NULL,
  issued_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME,
  revoked TINYINT(1) DEFAULT 0,
  ip_address VARCHAR(50),
  user_agent VARCHAR(255),
  FOREIGN KEY (user_id) REFERENCES app_user(user_id) ON DELETE CASCADE
);

-- Alerts
CREATE TABLE alert (
  id CHAR(36) PRIMARY KEY,
  type ENUM('anniversary','vacation','benefits','birthday') NOT NULL,
  title VARCHAR(255),
  message TEXT,
  employee_id CHAR(36),
  employee_name VARCHAR(255),
  alert_date DATETIME DEFAULT CURRENT_TIMESTAMP,
  priority ENUM('low','medium','high') DEFAULT 'low',
  is_read TINYINT(1) DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- KPI Summary
CREATE TABLE kpi_summary (
  id INT AUTO_INCREMENT PRIMARY KEY,
  snapshot_date DATE NOT NULL,
  total_earnings DECIMAL(18,2) DEFAULT 0,
  total_earnings_prev DECIMAL(18,2) DEFAULT 0,
  total_vacation_days INT DEFAULT 0,
  total_vacation_days_prev INT DEFAULT 0,
  average_benefits DECIMAL(12,2) DEFAULT 0,
  total_employees INT DEFAULT 0,
  active_alerts INT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Employee Lookup
CREATE TABLE employee_lookup (
  employee_id CHAR(36) PRIMARY KEY,
  hr_employee_db VARCHAR(50),
  payroll_employee_db VARCHAR(50),
  external_ref VARCHAR(255)
);

-- Indexes
CREATE INDEX idx_alert_employee ON alert(employee_id);
CREATE INDEX idx_user_email ON app_user(email);
