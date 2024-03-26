import { Node } from "reactflow";
import "./notebook.css";
const NotebookView = () => {
  const nodes: Node[] = [
    {
      id: "1", // required
      position: { x: 0, y: 0 }, // required
      data: {
        label:
          "\
        Lorem ipsum dolor sit amet. Ut sint velit non nihil sequi 33 fugit temporibus\
        in totam quod aut dolorem adipisci. Est maiores impedit qui recusandae nisi\
        non galisum dolor eum quod quisquam qui debitis esse qui debitis minus qui\
        voluptate similique. Nam temporibus quis qui impedit nesciunt qui ipsa fugiat\
        et deleniti repudiandae.\n\nEt quia excepturi sed incidunt fuga aut quod\
        odio est repellat sint non quaerat architecto aut magni earum! Ut velit odit\
        ut dignissimos laudantium ut illo mollitia in quaerat quisquam sed galisum\
        atque sit pariatur veniam ea error consectetur. Nam impedit commodi sed\
         laboriosam aperiam ea magni quibusdam.\n\nQui maiores laudantium\
        sed ipsa explicabo qui provident voluptates sit fuga accusantium. Cum autem \
        mollitia ex consequatur provident et soluta omnis ea culpa voluptas. In velit\
        blanditiis non quam enim et galisum omnis rem inventore rerum et galisum sunt\
        qui vitae natus. Ut voluptatem illum aut iste repellendus qui rerum maxime At\
        quis nulla.", // optional
      },
      style: {
        width: "400px",
      },
    },
  ];

  return (
    <div id="notebook-container">
      <div id="title">
        <h1>Notebook</h1>
      </div>
      <div id="workspace">
        <div
          id="flow-container"
          className="flex-item"
          style={{ width: "", height: "100%" }}
        ></div>
        <div id="notebook-tools" className="flex-item">
          {/* this will contain all of the files used in the current notebook */}
          <div id="resources" className="notebook-tool-item">
            the resources references in this notebook
          </div>
          <div id="notebook-chat" className="notebook-tool-item">
            some chat messages for the notebook
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotebookView;
