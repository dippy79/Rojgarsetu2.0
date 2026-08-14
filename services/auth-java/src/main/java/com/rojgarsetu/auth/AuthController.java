package com.rojgarsetu.auth;

import java.util.Map;
import java.util.Optional;

import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * Authentication Controller for user registration and login operations.
 * <p>Provides endpoints for user registration and authentication with JWT token generation.</p>
 */
@RestController
@RequestMapping("/auth")
public class AuthController {

    private final UserRepository userRepository;
    private final BCryptPasswordEncoder passwordEncoder;
    private final JwtService jwtService;

    /**
     * Constructor injection for dependency management.
     * <p>Ensures proper dependency injection and improves testability.</p>
     *
     * @param userRepository User data access layer
     * @param passwordEncoder Password encryption service
     * @param jwtService JWT token generation service
     */
    public AuthController(UserRepository userRepository, BCryptPasswordEncoder passwordEncoder, JwtService jwtService) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
        this.jwtService = jwtService;
    }

    /**
     * DTO Class to safely map JSON request body without casting issues.
     * <p>Provides type-safe request mapping for authentication endpoints.</p>
     */
    public static class AuthRequest {
        private String email;
        private String password;

        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
    }

    /**
     * User registration endpoint.
     * <p>Creates a new user account with hashed password and default CANDIDATE role.</p>
     *
     * @param body Authentication request containing email and password
     * @return ResponseEntity with success message or error details
     */
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

        String hashed = passwordEncoder.encode(password);
        User user = new User(email, hashed, "CANDIDATE");
        userRepository.save(user);

        return ResponseEntity.ok(Map.of(
            "message", "User registered successfully",
            "userId", String.valueOf(user.getId())
        ));
    }

    /**
     * User login endpoint.
     * <p>Authenticates user credentials and generates JWT token for session management.</p>
     *
     * @param body Authentication request containing email and password
     * @return ResponseEntity with JWT token and user role or error details
     */
    @PostMapping("/login")
    public ResponseEntity login(@RequestBody AuthRequest body) {
        String email = body.getEmail();
        String password = body.getPassword();

        if (email == null || password == null) {
            return ResponseEntity.badRequest().body(Map.of("error", "Email and password are required"));
        }

        Optional<User> userOpt = userRepository.findByEmail(email);

        if (userOpt.isPresent()) {
            User user = userOpt.get();
            if (passwordEncoder.matches(password, user.getPassword())) {
                String token = jwtService.generateToken(
                    user.getId().toString(),
                    user.getEmail(),
                    user.getRole() != null ? user.getRole() : "CANDIDATE"
                );
                return ResponseEntity.ok(Map.of(
                    "token", token,
                    "role", user.getRole() != null ? user.getRole() : "CANDIDATE"
                ));
            }
        }

        return ResponseEntity.status(401).body(Map.of("error", "Invalid credentials"));
    }
}