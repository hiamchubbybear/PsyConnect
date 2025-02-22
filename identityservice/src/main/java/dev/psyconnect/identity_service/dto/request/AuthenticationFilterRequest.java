package dev.psyconnect.identity_service.dto.request;

import java.util.Collection;
import java.util.List;
import java.util.stream.Collectors;

import dev.psyconnect.identity_service.model.Account;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class AuthenticationFilterRequest implements UserDetails {
    private String username;
    private String password;
    private List<GrantedAuthority> authorities;

    public AuthenticationFilterRequest(Account user) {
        this.username = user.getUsername();
        this.password = user.getPassword();
        this.authorities = user.getRole().stream()
                .map(role -> new SimpleGrantedAuthority("ROLE_" + role.getName()))
                .collect(Collectors.toList());
    }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return authorities;
    }

    @Override
    public boolean isAccountNonExpired() {
        return false;
    }

    @Override
    public boolean isAccountNonLocked() {
        return false;
    }

    @Override
    public boolean isCredentialsNonExpired() {
        return false;
    }

    @Override
    public boolean isEnabled() {
        return false;
    }
}
