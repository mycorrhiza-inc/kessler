import { SignIn } from "@clerk/nextjs";

export default function Page() {
  return (
    // <div className="flex content-center place-self-center">
      <SignIn path="/sign-in" signUpUrl="/sign-up" />
    // </div>
  );
}
