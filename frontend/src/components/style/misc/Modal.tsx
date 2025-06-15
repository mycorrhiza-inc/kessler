import { useEffect, useId } from "react";
const Modal = ({
  children,
  open,
  setOpen,
}: {
  children: React.ReactNode;
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  const modalID = useId();
  useEffect(() => {
    if (open) {
      // Confused as to why these are still necessary
      // @ts-ignore
      document.getElementById(modalID).showModal();
    } else {
      // Confused as to why these are still necessary
      // @ts-ignore
      document.getElementById(modalID).close();
    }
  }, [open]);
  return (
    <dialog id={modalID} className="modal ">
      <div
        className="modal-box bg-base-100 "
        style={{
          minHeight: "80vh",
          minWidth: "60vw",
        }}
      >
        {open && children}
      </div>
      <form
        method="dialog"
        className="modal-backdrop"
        onSubmit={() => setOpen(false)}
      >
        <button>close</button>
      </form>
    </dialog>
  );
};
export default Modal;

// Method for implemementing
// <form method="dialog">
//   {/* if there is a button in form, it will close the modal */}
//   <button
//     className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
//     onClick={() => {
//       setOpen(false);
//     }}
//   >
//     ✕
//   </button>
// </form>
