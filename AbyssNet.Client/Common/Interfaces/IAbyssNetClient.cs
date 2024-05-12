namespace AbyssNet.Client.Common.Interfaces;

public interface IAbyssNetClient
{
    Task InitializeClientAsync();
    Task PollAsync();
}