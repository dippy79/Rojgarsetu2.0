package com.rojgarsetu.auth;

import java.util.Map;
import java.util.Optional;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

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
        private String name;
        private String full_name;
        private String role;

        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
        public String getName() { return name != null ? name : full_name; }
        public void setName(String name) { this.name = name; }
        public void setFull_name(String full_name) { this.full_name = full_name; }
        public String getRole() { return role; }
        public void setRole(String role) { this.role = role; }
    }

    /**
     * User registration endpoint.
     * <p>Creates a new user account with hashed password and default CANDIDATE role.</p>
     *
     * @param body Authentication request containing email and password
     * @return ResponseEntity with success message or error details
     */
    @PostMapping("/register")
    public ResponseEntity<?> register(@RequestBody AuthRequest body) {
        String email = body.getEmail();
        String password = body.getPassword();

        if (email == null || !email.matches("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")) {
            return ResponseEntity.badRequest().body(Map.of("error", "A valid email is required"));
        }
        if (password == null || password.length() < 12) {
            return ResponseEntity.badRequest().body(Map.of("error", "Password must be at least 12 characters"));
        }

        if (userRepository.existsByEmail(email)) {
            return ResponseEntity.status(HttpStatus.CONFLICT).body(Map.of("error", "Email already exists"));
        }

        String name = body.getName();
        if (name == null || name.isEmpty()) {
            name = "User"; // Fallback to avoid NotNull constraint
        }

        String hashed = passwordEncoder.encode(password);
        User user = new User(name, email, hashed, body.getRole() != null ? body.getRole() : "CANDIDATE");
        userRepository.save(user);

        return ResponseEntity.ok(Map.of(
            "success", true,
            "message", "User registered successfully",
            "data", Map.of("user", userPayload(user))
        ));
    }

    /**
     * User login endpoint.
     * <p>Authenticates user credentials and generates JWT token for session management.</p>
     *
     * @param body Authentication request containing email and password
     * @param response HttpServletResponse to set HttpOnly cookie
     * @return ResponseEntity with JWT token and user role or error details
     */
    @PostMapping("/login")
    public ResponseEntity<?> login(@RequestBody AuthRequest body, HttpServletResponse response) {
        String email = body.getEmail();
        String password = body.getPassword();

        if (email == null || password == null || email.isBlank() || password.isBlank()) {
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

                // Set HttpOnly Cookie
                Cookie jwtCookie = new Cookie("access_token", token);
                jwtCookie.setHttpOnly(true);
                jwtCookie.setSecure(false); // Set to true in production with HTTPS
                jwtCookie.setPath("/");
                jwtCookie.setMaxAge(86400); // 24 hours
                response.addCookie(jwtCookie);

                return ResponseEntity.ok(Map.of(
                    "success", true,
                    "data", Map.of(
                        "user", userPayload(user)
                    )
                ));
            }
        }

        return ResponseEntity.status(401).body(Map.of(
            "success", false,
            "error", "Invalid credentials"
        ));
    }

    @GetMapping("/me")
    public ResponseEntity<?> currentUser(HttpServletRequest request) {
        Optional<User> user = authenticatedUser(request);
        if (user.isEmpty()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(Map.of("error", "Authentication required"));
        }
        return ResponseEntity.ok(Map.of("success", true, "data", userPayload(user.get())));
    }

    @PostMapping("/logout")
    public ResponseEntity<?> logout(HttpServletResponse response) {
        clearCookie(response, "access_token");
        clearCookie(response, "refresh_token");
        return ResponseEntity.ok(Map.of("success", true, "message", "Logged out successfully"));
    }

    private Optional<User> authenticatedUser(HttpServletRequest request) {
        if (request.getCookies() == null) return Optional.empty();
        for (Cookie cookie : request.getCookies()) {
            if ("access_token".equals(cookie.getName())) {
                try {
                    String userId = jwtService.validateToken(cookie.getValue()).getSubject();
                    return userRepository.findById(java.util.UUID.fromString(userId));
                } catch (RuntimeException ignored) {
                    return Optional.empty();
                }
            }
        }
        return Optional.empty();
    }

    private Map<String, String> userPayload(User user) {
        return Map.of(
            "id", user.getId().toString(),
            "email", user.getEmail(),
            "name", user.getName(),
            "role", user.getRole() != null ? user.getRole() : "CANDIDATE"
        );
    }

    private void clearCookie(HttpServletResponse response, String name) {
        Cookie cookie = new Cookie(name, "");
        cookie.setHttpOnly(true);
        cookie.setSecure(false);
        cookie.setPath("/");
        cookie.setMaxAge(0);
        response.addCookie(cookie);
    }
}
