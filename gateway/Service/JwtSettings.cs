namespace gateway.Service;

public class JwtSettings
{
    public string Issuer { get; set; } = string.Empty;
    public string[] Audiences { get; set; } = [];
    public bool RequireHttpsMetadata { get; set; }
    public string JwksUri { get; set; } = string.Empty;
}
