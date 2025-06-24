# 🚀 MediStock AI Backend

## Intelligent Medical Supply Chain Management

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://www.mysql.com/)

---

### **Overview**

The MediStock AI Backend is the robust engine powering an intelligent medical supply chain management system. Built with Go and backed by MySQL, it provides secure, scalable, and efficient APIs for managing users (hospitals, vendors, and internal staff), inventory, orders, and leverages machine learning for vendor scoring.

This repository contains the core API services, database interactions, and crucial background processes that ensure smooth operations and data-driven insights for healthcare supply chain optimization.

---
### **✨ Key Features**

* **User Management:** Secure registration, login, and profile management for various user types (Hospital, Vendor, Admin).
* **Role-Based Access Control (RBAC):** Granular authorization ensures users only access resources relevant to their roles.
* **JWT Authentication:** Secure API access using Access and Refresh Tokens for seamless session management across multiple devices.
* **Multi-Hospital Support:** Designed to handle multiple distinct hospital entities.
* **Vendor Management:** APIs for vendor information, supply price listings, and performance tracking.
* **Inventory Management:** Track medical supplies for each hospital, including stock levels, reorder points, and expiry dates.
* **Order Management:** Streamlined process for hospitals to place orders with vendors, track order status, and manage order items.
* **ML-Powered Vendor Scoring:** A robust background process that periodically fetches vendor data, sends it to an external ML service for worthiness scoring, and updates the scores in the database for informed decision-making.
* **MySQL Database:** Reliable and efficient data storage for all application entities.

---

### **🛠️ Technologies Used**

* **Go (Golang):** Core backend language.
* **MySQL:** Relational database for persistent storage.
* **`gorilla/mux` or `gin-gonic/gin`:** (Choose one and specify) High-performance HTTP router for API endpoints.
* **`jmoiron/sqlx`:** (Optional but recommended for `database/sql`) A powerful extension to Go's `database/sql` for easier mapping.
* **`golang.org/x/crypto/bcrypt`:** For secure password hashing.
* **`github.com/golang-jwt/jwt/v5`:** For JSON Web Token (JWT) handling.
* **`github.com/robfig/cron/v3` or `github.com/go-co-op/gocron/v2`:** (Choose one and specify) For scheduling background tasks.
* **`github.com/joho/godotenv`:** (Optional) For loading environment variables from `.env` files.
* **`uuid` package:** For generating UUIDs for IDs.

---

### **🔐 Authentication & Authorization**

This backend implements a robust JWT-based authentication system:

* Access Tokens: Short-lived tokens sent in the Authorization: Bearer header for every protected API request.
* Refresh Tokens: Long-lived tokens (stored securely in HttpOnly cookies) used to obtain new access tokens when the current one expires, ensuring seamless user sessions.
* Middleware: Middleware is used to validate tokens and enforce role-based access control for all protected routes.

---

### **⚙️ Background Processes**

A critical background process runs periodically (configured for every 24 hours) to:

* Fetch vendor data in batches.
* Send this data to an external Machine Learning (ML) service.
* Receive a ml_worth_score from the ML service.
* Update the ml_worth_score in the database for the respective vendors.
* This ensures the vendor scoring is continuously updated without impacting real-time API performance.

---

### **🚀 Getting Started**

Follow these steps to get the MediStock AI Backend up and running on your local machine.

#### **Prerequisites**

* **Go:** [Go 1.22+](https://golang.org/doc/install)
* **MySQL:** [MySQL Server 8.0+](https://dev.mysql.com/doc/refman/8.0/en/installing.html)
* **Git:** [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
