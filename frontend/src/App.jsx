import "./App.css";
import {
  SelectedDirectory,
  StartHTTPServer,
  StopHTTPServer,
  SetDraftRootPath,
  AutoDetectDraftRootPath,
} from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from "../wailsjs/runtime/runtime";
import Swal from "sweetalert2";
import { useState, useEffect, useRef } from "react";

function App() {
  // Drafts folder state and update function
  const [draftsFolder, setDraftsFolder] = useState(
    localStorage.getItem("draftsFolder") || ""
  );
  const updateDraftsFolder = (e) => {
    setDraftsFolder(e.target.value);
  };
  useEffect(() => {
    localStorage.setItem("draftsFolder", draftsFolder);
    SetDraftRootPath(draftsFolder).then(() => {
      console.log("Drafts folder updated successfully");
    }).catch((err) => {
      Swal.fire({
        title: "Oops...",
        text: "Error updating drafts folder: " + err,
        icon: "error",
      });
    });
  }, [draftsFolder])

  // Open the directory dialog box
  const selectedDirectoryHandle = () => {

    // 自动检测草稿目录
    AutoDetectDraftRootPath().then((path) => {
      // console.log(path);
      // setDraftsFolder(path);
      // Swal.fire({
      //   title: "Drafts folder detected",
      //   text: "Drafts folder detected: " + path,
      //   icon: "success",
      // });

      // 提示用户已经自动检测到地址，是否直接使用
      Swal.fire({
        title: "Drafts folder detected",
        text: "Drafts folder detected: " + path,
        icon: "success",
        showCancelButton: true,
        confirmButtonText: "Yes, use this folder",
        cancelButtonText: "No, select another folder",
      }).then((result) => {
        if (result.isConfirmed) {
          setDraftsFolder(path);
        } else {
          SelectedDirectory()
            .then((path) => {
              console.log(path);
              setDraftsFolder(path);
            })
            .catch((err) => {
              Swal.fire({
                title: "Oops...",
                text: "Error opening directory dialog box:" + err,
                icon: "error",
              });
            });
        }
      });
    }).catch((err) => {
      SelectedDirectory()
        .then((path) => {
          console.log(path);
          setDraftsFolder(path);
        })
        .catch((err) => {
          Swal.fire({
            title: "Oops...",
            text: "Error opening directory dialog box:" + err,
            icon: "error",
          });
        });
    });


  };

  const [runningState, setRunningState] = useState(
    localStorage.getItem("runningState") === "true" || false
  );
  const updateRunning = () => {
    setRunningState(true);
    localStorage.setItem("runningState", true);
  };
  const updateStopped = () => {
    setRunningState(false);
    localStorage.setItem("runningState", false);
  };
  useEffect(() => {
    const fetchRunningState = async () => {
      try {
        const response = await fetch("http://localhost:9507");
        await response.text();
        updateRunning();
      } catch (error) {
        console.log(error);
        updateStopped();
      }
    };

    fetchRunningState();
  }, []);

  // Start the HTTP server
  const startHandle = () => {
    StartHTTPServer()
      .then(() => {
        Swal.fire({
          title: "Success",
          text: "HTTP server started successfully",
          icon: "success",
          timer: 2000,
          timerProgressBar: true,
          showConfirmButton: false,
        });
        updateRunning();
      })
      .catch((err) => {
        Swal.fire({
          title: "Oops...",
          text: "Error starting HTTP server: " + err,
          icon: "error",
        });
      });
  };

  // Stop the HTTP server
  const stopHandle = () => {
    StopHTTPServer()
      .then(() => {
        Swal.fire({
          title: "Success",
          text: "HTTP server stopped successfully",
          icon: "success",
          timer: 2000,
          timerProgressBar: true,
          showConfirmButton: false,
        });
        updateStopped();
      })
      .catch((err) => {
        Swal.fire({
          title: "Oops...",
          text: "Error stopping HTTP server: " + err,
          icon: "error",
        });
      });
  };

  // Log event listener
  const [logs, setLogs] = useState([]);
  useEffect(() => {
    EventsOn("logs", function (data) {
      // data {"type": "log", "message": "Hello, world!"}
      data = JSON.parse(data);

      setLogs((prev) => [...prev, data.message]);
    });
    return () => {
      EventsOff("logs");
    };
  }, []);

  // Scroll to the bottom of the log container
  const logContainerRef = useRef(null);
  useEffect(() => {
    if (logContainerRef.current) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight;
    }
  }, [logs]);

  return (
    <div id="App" className="p-5">
      <div className="flex flex-col gap-5">
        <label className="form-control w-full">
          <div className="label">
            <span className="label-text">Drafts Folder</span>
          </div>

          <label className="input input-bordered flex items-center gap-2">
            <input
              type="text"
              className="grow"
              placeholder="Enter or select your drafts folder"
              value={draftsFolder}
              onChange={updateDraftsFolder}
            />
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 1024 1024"
              height="2em"
              width="2em"
              onClick={selectedDirectoryHandle}
            >
              <path fill="#D9D9D9" d="M159 768h612.3l103.4-256H262.3z" />
              <path d="M928 444H820V330.4c0-17.7-14.3-32-32-32H473L355.7 186.2a8.15 8.15 0 0 0-5.5-2.2H96c-17.7 0-32 14.3-32 32v592c0 17.7 14.3 32 32 32h698c13 0 24.8-7.9 29.7-20l134-332c1.5-3.8 2.3-7.9 2.3-12 0-17.7-14.3-32-32-32zM136 256h188.5l119.6 114.4H748V444H238c-13 0-24.8 7.9-29.7 20L136 643.2V256zm635.3 512H159l103.3-256h612.4L771.3 768z" />
            </svg>
          </label>
        </label>

        {!runningState ? (
          <button className="btn w-full" onClick={startHandle}>
            Start
          </button>
        ) : (
          <button className="btn w-full" onClick={stopHandle}>
            Stop
          </button>
        )}

        {runningState ? (
          <div className="flex gap-2 justify-center items-center font-bold">
            <div className="badge badge-success badge-sm"></div>Running
          </div>
        ) : (
          <div className="flex gap-2 justify-center items-center font-bold">
            <div className="badge badge-error badge-sm"></div>Stopped
          </div>
        )}
      </div>
      <div
        className="mockup-code mt-10 h-80 overflow-y-auto text-xs"
        ref={logContainerRef}
      >
        {logs.map((log, index) => (
          <pre data-prefix=">" key={index}>
            <code>{log}</code>
          </pre>
        ))}
      </div>

      <p className="mt-8 text-center text-xs text-gray-400">Copyright © 2024 - All right reserved by <a href="https://x.com/iHunterDev" target="_blank">@iHunterDev</a></p>
      <p className="mt-1 text-center text-xs text-gray-400">Version: v0.1.2</p>
    </div>
  );
}

export default App;
