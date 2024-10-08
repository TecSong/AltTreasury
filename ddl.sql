CREATE TABLE `withdraw_claims` (
  `id` int NOT NULL AUTO_INCREMENT,
  `staff_id` int NOT NULL,
  `amount` decimal(18,8) NOT NULL,
  `token_address` varchar(42) NOT NULL,
  `recipient_address` varchar(42) NOT NULL,
  `status` enum('pending','approved','rejected','executed') NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_withdrawal_claim_status` (`status`),
  KEY `idx_withdrawal_claim_staff_id` (`staff_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `withdraw_claim_confirmations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `withdraw_claim_id` int DEFAULT NULL,
  `manager_id` int NOT NULL,
  `action_type` enum('approve','reject') NOT NULL,
  `confirmed_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_claim_manager_action` (`withdraw_claim_id`, `manager_id`, `action_type`),
  KEY `idx_withdrawal_claim_confirmation_claim_id` (`withdraw_claim_id`),
  KEY `idx_withdrawal_claim_confirmation_confirmed_at` (`confirmed_at`),
  CONSTRAINT `withdraw_claim_confirmations_ibfk_1` FOREIGN KEY (`withdraw_claim_id`) REFERENCES `withdraw_claims` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE withdrawal_claim_confirmation (
    id INT AUTO_INCREMENT PRIMARY KEY,
    withdrawal_claim_id INTEGER NOT NULL,  // 注意这里是 withdrawal_claim_id，不是 withdraw_claim_id
    manager_id INTEGER NOT NULL,
    action_type ENUM('approve', 'reject') NOT NULL,
    confirmed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (withdrawal_claim_id) REFERENCES withdrawal_claim(id),
    UNIQUE KEY unique_confirmation (withdrawal_claim_id, manager_id, action_type)
);
