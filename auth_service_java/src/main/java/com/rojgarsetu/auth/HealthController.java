package com.rojgarsetu.auth;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.sql.DataSource;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/")
public class HealthController {

    @Autowired
    private DataSource dataSource;

    @GetMapping("health")
    public ResponseEntity<Map<String, Object>> health() {
        Map<String, Object> response = new HashMap<>();
        response.put("status", "UP");
        response.put("service", "auth-service");
        response.put("version", "2.0.0");
        
        // Check database connection
        try {
            dataSource.getConnection().close();
            response.put("database", "UP");
        } catch (Exception e) {
            response.put("database", "DOWN");
            response.put("dbError", e.getMessage());
        }
        
        return ResponseEntity.ok(response);
    }

    @GetMapping("actuator/health")
    public ResponseEntity<Map<String, Object>> actuatorHealth() {
        Map<String, Object> response = new HashMap<>();
        response.put("status", "UP");
        return ResponseEntity.ok(response);
    }
}

