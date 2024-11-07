"use client";
import { Suspense, useState } from "react";
import LoadingSpinner from "../styled-components/LoadingSpinner";

const SuspenseTest = () => {
  const [count, setCount] = useState(0);
  return (
    <div>
      <button
        className="btn btn-primary"
        onClick={() => setCount((prev) => prev + 1)}
      >
        This is a test button
      </button>
      <Suspense fallback={<LoadingSpinner />}>
        <SuspenseResult num={count} />
      </Suspense>
    </div>
  );
};

const SuspenseResult = async ({ num }: { num: number }) => {
  async function fetchSomeResults(num: number) {
    await new Promise((r) => setTimeout(r, 2000));
    return Array.from({ length: num }, () => Math.random());
  }
  const results = await fetchSomeResults(num);
  return (
    <table className="table table-zebra">
      {/* head */}
      <thead>
        <tr>
          <th>Result ID</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        {results.map((val, index) => (
          <tr>
            <td>{index}</td>
            <td>{val}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default SuspenseTest;
