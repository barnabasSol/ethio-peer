using Microsoft.IdentityModel.Tokens;

namespace gateway.Service;

public interface IJwksProvider
{
    IReadOnlyCollection<SecurityKey> GetSigningKeys();
    Task RefreshKeysAsync();
}

public class JwksProvider(string jwksUrl) : IJwksProvider
{
    private readonly HttpClient _httpClient = new();

    private readonly string _jwksUrl = jwksUrl;
    private List<SecurityKey> _keys = [];

    public IReadOnlyCollection<SecurityKey> GetSigningKeys() => _keys;

    public async Task RefreshKeysAsync()
    {
        var jwksJson = await _httpClient.GetStringAsync(_jwksUrl);
        var jwks = new JsonWebKeySet(jwksJson);
        _keys = [.. jwks.GetSigningKeys()];
    }
}
