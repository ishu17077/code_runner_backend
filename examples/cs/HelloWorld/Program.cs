using System;
public class Program
{
    public static void Main(string[] args)
    {
        // Console.Write("Please enter an integer: ");
        string input = Console.ReadLine(); int number;
        bool success = int.TryParse(input, out number);
        if (success)
        {
            if (number % 2 == 0)
            {
                Console.WriteLine("Yes");
            }
            else
            {
                Console.WriteLine("No");
            }
        }
        else
        {
            Console.WriteLine("Invalid input. Please enter a valid integer.");
        }
    }
}
