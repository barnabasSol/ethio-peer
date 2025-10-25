using System.Text;
using System.Text.Json;
using Microsoft.Extensions.Configuration;

public class GeminiCaller
{
    private readonly IConfiguration _IConfiguration;
    public GeminiCaller(IConfiguration configuration)
    {
        _IConfiguration = configuration;
    }
    private readonly string _systemInstruction = "Act as a Senior Academic Analyst. Your task is to analyze raw weekly discussion transcripts and generate a three-part executive summary for a Director. Focus on trends and pain points. Return the output in Markdown. Don't mention to and from.Try to focus on academic discussion than students random chats.";
    private readonly string _userQuery = "Based on the raw discussion logs provided below, please identify the top 3 most vibrant rooms, the main discussed topics(points) in those rooms, and the most critical academic insights or recurring gaps that require faculty attention.\nLogs: ";


    public async Task<string> CallGeminiApiAsync(string discussionLogs)
    {
        var client = new HttpClient();
        client.DefaultRequestHeaders.Add("x-goog-api-key", _IConfiguration["Gemini:ApiKey"]!);
        var requestBody = new
        {
            contents = new[] {
        new { parts = new[] { new { text = _userQuery+discussionLogs } } }
    },
            systemInstruction = new
            {
                parts = new[]
    {
        new{text=_systemInstruction}
    }
            }
        };
        var response = await client.PostAsync(
            _IConfiguration["Gemini:Baseurl"],
            new StringContent(JsonSerializer.Serialize(requestBody), Encoding.UTF8, "application/json")
        );
        string result = await response.Content.ReadAsStringAsync();
        
        try
        {
            var parsedJson = JsonDocument.Parse(result);

            if (parsedJson.RootElement.TryGetProperty("candidates", out var candidates) &&
                candidates.GetArrayLength() > 0 &&
                candidates[0].TryGetProperty("content", out var content) &&
                content.TryGetProperty("parts", out var parts) &&
                parts.GetArrayLength() > 0 &&
                parts[0].TryGetProperty("text", out var textElement))
            {
                return textElement.GetString() ?? "Gemini returned empty text";
            }
            else
            {
                return "Gemini response missing expected structure";
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine("Raw Gemini response: " + result);
            return $"Error parsing Gemini response: {ex.Message}";
        }
    }


}