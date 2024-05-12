namespace AbyssNet.Client.Common;
public class AbyssNetClientWebSocketOptions
{
    // The URI of the server to connect to.
    public required string ServerURI { get; set; } = string.Empty;
    
    // The interval to keep the connection alive. Default is 10 seconds.
    public TimeSpan KeepAliveInterval { get; set; } = TimeSpan.FromSeconds(10);

    // The timeout for the connection. Default is 120 seconds.
    public TimeSpan ConnectionTimeout { get; set; } = TimeSpan.FromSeconds(120);
}