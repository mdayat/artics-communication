import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/Button";
import { Card, CardContent, CardFooter } from "@/components/ui/Card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/Form";
import { Input } from "@/components/ui/Input";
import { Link, useNavigate } from "react-router";
import { axiosInstance } from "@/lib/axios";
import { toast } from "sonner";
import type { RegisterRequest, UserResponse } from "@/dtos/user";
import type { AxiosResponse } from "axios";

const formSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
  email: z.string().email({
    message: "Please enter a valid email address.",
  }),
  password: z.string().min(8, {
    message: "Password must be at least 8 characters.",
  }),
});

function Registration() {
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setIsLoading(true);
    try {
      const res = await axiosInstance.post<
        UserResponse,
        AxiosResponse<UserResponse>,
        RegisterRequest
      >("/auth/register", {
        username: values.username,
        email: values.email,
        password: values.password,
      });

      if (res.status === 201) {
        toast.success("Registration success, please login to continue", {
          richColors: true,
        });
        navigate("/login", { replace: true });
      } else if (res.status === 400) {
        toast.error("Please check your username, email, and password again", {
          richColors: true,
        });
      } else if (res.status === 409) {
        toast.error("Email already registered", {
          richColors: true,
        });
      } else {
        throw new Error(`unknown status code: ${res.status}`);
      }
    } catch (error) {
      console.error(new Error("failed to register", { cause: error }));
      toast.error("Registration failed, please try again", {
        richColors: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-4 md:p-24">
      <div className="w-full max-w-md">
        <h1 className="mb-6 text-center text-3xl font-bold tracking-tight">
          Create an account
        </h1>

        <Card className="border-slate-200 shadow-sm">
          <CardContent className="pt-6">
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                autoComplete="off"
                className="space-y-6"
              >
                <FormField
                  control={form.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Username</FormLabel>
                      <FormControl>
                        <Input
                          disabled={isLoading}
                          placeholder="John"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        This will be your public display name.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input
                          disabled={isLoading}
                          type="email"
                          placeholder="john@example.com"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        We'll never share your email with anyone else.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Password</FormLabel>
                      <FormControl>
                        <Input
                          disabled={isLoading}
                          type="password"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Must be at least 8 characters long.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" className="w-full" disabled={isLoading}>
                  {isLoading ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Creating account...
                    </>
                  ) : (
                    "Register"
                  )}
                </Button>
              </form>
            </Form>
          </CardContent>

          <CardFooter className="flex justify-center border-t p-4">
            <p className="text-sm text-slate-500">
              Already have an account?{" "}
              <Link
                to="/login"
                className="text-slate-900 underline hover:text-slate-700"
              >
                Sign in
              </Link>
            </p>
          </CardFooter>
        </Card>
      </div>
    </main>
  );
}

export { Registration };
