using System;
using System.IO;
using System.Reflection;

public class Reflect
{
    public static void GetAllClassesAndMethodsOfAssembly(string name)
    {
        Assembly assem1 = Assembly.LoadFrom(name);
        //Another Way
        //Assembly assem2 = Assembly.Load(name);
        //Get List of Class Name
        Console.WriteLine("Assembly FullName : " + assem1.FullName);
        Console.WriteLine("Assembly EntryPoint : " + assem1.EntryPoint);

        Type[] types = assem1.GetTypes();
        foreach (Type tc in types)
        {
    
            if (tc.IsAbstract)
            {
                Console.WriteLine("Abstract Class : " + tc.Name);
            }
            else if (tc.IsPublic)
            {
                Console.WriteLine("Public Class : " + tc.Name);
            }
            else if (tc.IsSealed)
            {
                Console.WriteLine("Sealed Class : " + tc.Name);
            }  

            Console.WriteLine("Type.Namespace : " + tc.Namespace);
            Console.WriteLine("Type.FullName : " + tc.FullName);
            Console.WriteLine("Type.AssemblyQualifiedName : " + tc.AssemblyQualifiedName);

            //Get List of Method Names of Class
            MemberInfo[] methodName = tc.GetMethods();

            foreach (MemberInfo method in methodName)
            {
                if (method.ReflectedType.IsPublic)
                {
                    Console.WriteLine("Public Method : " + method.Name.ToString());
                }
                else
                {
                    Console.WriteLine("Non-Public Method : " + method.Name.ToString());
                }
            }
        }
    }

    static void Main(string[] args)
    {
        Console.WriteLine("Input File: " + args[0]);
        var filename = args[0];
        Console.WriteLine(File.Exists(filename) ? "File exists." : "File does not exist.");
        GetAllClassesAndMethodsOfAssembly(filename);
    }  
}