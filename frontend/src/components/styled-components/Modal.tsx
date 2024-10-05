import { useEffect, useState } from "react";
const Modal = ({
  children,
  open,
  setOpen,
}: {
  children: React.ReactNode;
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  const uuid = (Math.random() + 1).toString(36).substring(7);
  useEffect(() => {
    if (open) {
      // Oh come on, the docid is always not null, its defined right below you
      // @ts-ignore
      document.getElementById(`doc_modal_${uuid}`).showModal();
    } else {
      // Oh come on, the docid is always not null, its defined right below you
      // @ts-ignore
      document.getElementById(`doc_modal_${uuid}`).close();
    }
  }, [open]);
  return (
    <dialog id={`doc_modal_${uuid}`} className="modal ">
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
