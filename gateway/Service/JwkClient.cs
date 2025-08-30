namespace gateway.Service;

public static class JwkClient
{
    public static async Task<string> GetJWKSFromHttp(string url)
    {
        using var httpClient = new HttpClient();
        var response = await httpClient.GetAsync(url);

        if (!response.IsSuccessStatusCode)
        {
            throw new Exception($"Unable to retrieve JWKS from {url}: {response.StatusCode}");
        }

        return await response.Content.ReadAsStringAsync();
    }
}
