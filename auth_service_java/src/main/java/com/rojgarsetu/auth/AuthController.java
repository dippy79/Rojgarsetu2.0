package com.rojgarsetu.auth;

import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.web.bind.annotation.*;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.*;

@RestController
@RequestMapping("/api/auth")
public class AuthController {

    @Value("${jwt.secret:your-super-secret-jwt-key-change-in-production}")
    private String jwtSecret;

    @Value("${jwt.expiration:86400000}")
    private long jwtExpiration;

    // In-memory user store (for demo - should use database in production)
    private final Map<String, User> users = new HashMap<>();
    private final BCryptPasswordEncoder passwordEncoder = new BCryptPasswordEncoder();

    public AuthController() {
        // Create a demo user
        User demoUser = new User();
        demoUser.id = UUID.randomUUID().toString();
        demoUser.email = "demo@rojgarsetu.com";
        demoUser.passwordHash = passwordEncoder.encode("demo123");
        demoUser.name = "Demo User";
        users.put(demoUser.email, demoUser);
    }

    @PostMapping("/register")
    public ResponseEntity<Map<String, Object>> register(@RequestBody RegisterRequest request) {
        Map<String, Object> response = new HashMap<>();

        if (users.containsKey(request.getEmail())) {
            response.put("status", "error");
            response.put("message", "Email already registered");
            return ResponseEntity.status(HttpStatus.CONFLICT).body(response);
        }

        User newUser = new User();
        newUser.id = UUID.randomUUID().toString();
        newUser.email = request.getEmail();
        newUser.passwordHash = passwordEncoder.encode(request.getPassword());
        newUser.name = request.getName() != null ? request.getName() : "";
        users.put(newUser.email, newUser);

        response.put("status", "success");
        response.put("message", "User registered successfully");
        response.put("userId", newUser.id);

        return ResponseEntity.ok(response);
    }

    @PostMapping("/login")
    public ResponseEntity<Map<String, Object>> login(@RequestBody LoginRequest request) {
        Map<String, Object> response = new HashMap<>();

        User user = users.get(request.getEmail());
        if (user == null || !passwordEncoder.matches(request.getPassword(), user.passwordHash)) {
            response.put("status", "error");
            response.put("message", "Invalid email or password");
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
        }

        // Generate JWT token
        SecretKey key = Keys.hmacShaKeyFor(jwtSecret.getBytes(StandardCharsets.UTF_8));
        String token = Jwts.builder()
                .setSubject(user.id)
                .claim("email", user.email)
                .claim("name", user.name)
                .setIssuedAt(new Date())
                .setExpiration(new Date(System.currentTimeMillis() + jwtExpiration))
                .signWith(key, SignatureAlgorithm.HS256)
                .compact();

        response.put("status", "success");
        response.put("token", token);
        response.put("user", Map.of(
                "id", user.id,
                "email", user.email,
                "name", user.name
        ));

        return ResponseEntity.ok(response);
    }

    @PostMapping("/validate")
    public ResponseEntity<Map<String, Object>> validate(@RequestHeader("Authorization") String authHeader) {
        Map<String, Object> response = new HashMap<>();

        try {
            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                response.put("valid", false);
                return ResponseEntity.ok(response);
            }

            String token = authHeader.substring(7);
            SecretKey key = Keys.hmacShaKeyFor(jwtSecret.getBytes(StandardCharsets.UTF_8));
            
            var claims = Jwts.parserBuilder()
                    .setSigningKey(key)
                    .build()
                    .parseClaimsJws(token);

            response.put("valid", true);
            response.put("userId", claims.getBody().getSubject());
            response.put("email", claims.getBody().get("email"));

            return ResponseEntity.ok(response);
        } catch (Exception e) {
            response.put("valid", false);
            response.put("error", e.getMessage());
            return ResponseEntity.ok(response);
        }
    }

    @GetMapping("/me")
    public ResponseEntity<Map<String, Object>> getCurrentUser(@RequestHeader("Authorization") String authHeader) {
        Map<String, Object> response = new HashMap<>();

        try {
            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                response.put("status", "error");
                response.put("message", "No token provided");
                return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
            }

            String token = authHeader.substring(7);
            SecretKey key = Keys.hmacShaKeyFor(jwtSecret.getBytes(StandardCharsets.UTF_8));
            
            var claims = Jwts.parserBuilder()
                    .setSigningKey(key)
                    .build()
                    .parseClaimsJws(token);

            String userId = claims.getBody().getSubject();
            
            // Find user by ID
            User user = users.values().stream()
                    .filter(u -> u.id.equals(userId))
                    .findFirst()
                    .orElse(null);

            if (user == null) {
                response.put("status", "error");
                response.put("message", "User not found");
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body(response);
            }

            response.put("status", "success");
            response.put("user", Map.of(
                    "id", user.id,
                    "email", user.email,
                    "name", user.name
            ));

            return ResponseEntity.ok(response);
        } catch (Exception e) {
            response.put("status", "error");
            response.put("message", "Invalid token");
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(response);
        }
    }

    // Inner classes
    static class User {
        String id;
        String email;
        String passwordHash;
        String name;
    }

    static class RegisterRequest {
        private String email;
        private String password;
        private String name;

        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
        public String getName() { return name; }
        public void setName(String name) { this.name = name; }
    }

    static class LoginRequest {
        private String email;
        private String password;

        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
    }
}

