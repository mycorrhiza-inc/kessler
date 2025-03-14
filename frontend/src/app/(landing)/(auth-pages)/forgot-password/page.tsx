import { forgotPasswordAction } from "@/app/actions";
import {
  FormMessage,
  Message,
} from "@/components/supabasetutorial/form-message";
import { SubmitButton } from "@/components/supabasetutorial/submit-button";
import { Input } from "@/components/supabasetutorial/ui/input";
import { Label } from "@/components/supabasetutorial/ui/label";

export default function ForgotPassword({
  searchParams,
}: {
  searchParams: Message;
}) {
  return (
    <>
      <form className="flex-1 flex flex-col w-full gap-2 text-base-content [&>input]:mb-6 min-w-64 max-w-64 mx-auto">
        <div>
          <h1 className="text-2xl font-medium">Reset Password</h1>
          <p className="text-sm text-base-content">
            Already have an account?{" "}
            <a className="underline" href="/sign-in">
              Sign in
            </a>
          </p>
        </div>
        <div className="flex flex-col gap-2 [&>input]:mb-3 mt-8">
          <Label htmlFor="email">Email</Label>
          <Input name="email" placeholder="you@example.com" required />
          <SubmitButton formAction={forgotPasswordAction}>
            Reset Password
          </SubmitButton>
          <FormMessage message={searchParams} />
        </div>
      </form>
    </>
  );
}
