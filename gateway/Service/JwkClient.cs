namespace gateway.Service;

public static class JwkClient
{
    public static async Task<string> GetJWKSFromHttp(
        string url,
        int maxRetries = 5,
        int baseDelayMs = 500
    )
    {
        using var httpClient = new HttpClient();
        var random = new Random();

        for (int attempt = 0; attempt < maxRetries; attempt++)
        {
            try
            {
                var response = await httpClient.GetAsync(url);
                if (response.IsSuccessStatusCode)
                {
                    return await response.Content.ReadAsStringAsync();
                }
                else
                {
                    throw new Exception(
                        $"Unable to retrieve JWKS from {url}: {response.StatusCode}"
                    );
                }
            }
            catch (Exception ex) when (attempt < maxRetries - 1)
            {
                Console.WriteLine(ex.Message);
                int delay = (int)(baseDelayMs * Math.Pow(2, attempt));
                delay += random.Next(0, 100);
                await Task.Delay(delay);
            }
        }
        var finalResponse = await httpClient.GetAsync(url);
        if (!finalResponse.IsSuccessStatusCode)
        {
            throw new Exception($"Unable to retrieve JWKS from {url}: {finalResponse.StatusCode}");
        }
        return await finalResponse.Content.ReadAsStringAsync();
    }
}
