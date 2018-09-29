using System;
using System.IO;
using System.Reflection;

public class Class1
{
    public static void GetAllClassesAndMethodsOfAssembly(string name)
    {
        Assembly assem1 = Assembly.Load(AssemblyName.GetAssemblyName(name));
        //Another Way
        Assembly assem2 = Assembly.Load(name);
        //Get List of Class Name

        Type[] types = assem1.GetTypes();
        foreach(Type tc in types)
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
        GetAllClassesAndMethodsOfAssembly(args[0]);
    }  
}