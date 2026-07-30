package com.rojgarsetu.auth;

import java.util.Map;
import java.util.Optional;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/auth")
public class AuthController {

    @Autowired
    private UserRepository userRepository;

    // DTO Class to safely map JSON request body without casting issues
    static class AuthRequest {
        private String email;
        private String password;

        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
    }

    @PostMapping("/register")
    public ResponseEntity register(@RequestBody AuthRequest body) {
        String email = body.getEmail();
        String password = body.getPassword();

        if (email == null || password == null) {
            return ResponseEntity.badRequest().body(Map.of("error", "Email and password are required"));
        }

        if (userRepository.existsByEmail(email)) {
            return ResponseEntity.badRequest().body(Map.of("error", "Email already exists"));
        }

        User user = new User(email, password, "CANDIDATE");
        userRepository.save(user);

        return ResponseEntity.ok(Map.of(
            "message", "User registered successfully", 
            "userId", String.valueOf(user.getId())
        ));
    }

    @PostMapping("/login")
    public ResponseEntity login(@RequestBody AuthRequest body) {
        String email = body.getEmail();
        String password = body.getPassword();

        if (email == null || password == null) {
            return ResponseEntity.badRequest().body(Map.of("error", "Email and password are required"));
        }

        Optional<User> userOpt = userRepository.findByEmail(email);
        
        if (userOpt.isPresent()) {
            User user = userOpt.get(); // Explicitly returns User, not Object
            if (user.getPassword().equals(password)) {
                return ResponseEntity.ok(Map.of(
                    "token", "dummy-jwt-token-" + user.getId(),
                    "role", user.getRole() != null ? user.getRole() : "CANDIDATE"
                ));
            }
        }

        return ResponseEntity.status(401).body(Map.of("error", "Invalid credentials"));
    }
}