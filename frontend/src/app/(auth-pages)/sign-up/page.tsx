import { signUpAction } from "@/app/actions";
import {
  FormMessage,
  Message,
} from "@/components/supabasetutorial/form-message";
import { SubmitButton } from "@/components/supabasetutorial/submit-button";
import { Input } from "@/components/supabasetutorial/ui/input";
import { Label } from "@/components/supabasetutorial/ui/label";

export default function Signup({ searchParams }: { searchParams: Message }) {
  if ("message" in searchParams) {
    return (
      <div className="w-full flex-1 flex items-center h-screen sm:max-w-md justify-center gap-2 p-4">
        <FormMessage message={searchParams} />
      </div>
    );
  }

  return (
    <>
      <form className="flex flex-col min-w-64 max-w-64 mx-auto">
        <h1 className="text-2xl font-medium">Sign up</h1>
        <p className="text-sm text-base-content">
          Already have an account?{" "}
          <a className="font-medium underline" href="/sign-in">
            Sign in
          </a>
        </p>
        <div className="flex flex-col gap-2 [&>input]:mb-3 mt-8">
          <Label htmlFor="email">Email</Label>
          <Input name="email" placeholder="you@example.com" required />
          <Label htmlFor="password">Password</Label>
          <Input
            type="password"
            name="password"
            placeholder="Your password"
            minLength={6}
            required
          />
          <SubmitButton formAction={signUpAction} pendingText="Signing up...">
            Sign up
          </SubmitButton>
          <FormMessage message={searchParams} />
        </div>
      </form>
    </>
  );
}
