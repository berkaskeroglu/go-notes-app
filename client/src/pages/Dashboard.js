import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import API from "../api";

import { faEdit, faSave, faTrash } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function Dashboard() {
    const [notes, setNotes] = useState([]);
    const [title, setTitle] = useState("");
    const [body, setBody] = useState("");
    const [editNoteId, setEditNoteId] = useState(null);
    const [editTitle, setEditTitle] = useState("");
    const [editBody, setEditBody] = useState("");
    const navigate = useNavigate();

    const fetchNotes = async () => {
        try {
            const res = await API.get("/notes");
            if (Array.isArray(res.data)) {
                setNotes(res.data);
            } else {
                setNotes([]);
            }
        } catch (err) {
            console.error("Note fetch error:", err);
            setNotes([]);
            navigate("/login");
        }
    };

    const addNote = async () => {
        if (!title || !body) return;
        await API.post("/notes", { title, body });
        setTitle("");
        setBody("");
        fetchNotes();
    };

    const deleteNote = async (id) => {
        await API.delete(`/notes/${id}`);
        fetchNotes();
    };

    const startEdit = (note) => {
        setEditNoteId(note.id);
        setEditTitle(note.title);
        setEditBody(note.body);
    };

    const saveEdit = async () => {
        await API.put(`/notes/${editNoteId}`, {
            title: editTitle,
            body: editBody,
        });
        setEditNoteId(null);
        setEditTitle("");
        setEditBody("");
        fetchNotes();
    };


    const logout = () => {
        localStorage.removeItem("token");
        navigate("/login");
    };

    useEffect(() => {
        fetchNotes();
    }, []);

    return (
        <div className="container">
            <h2 className="center">Dashboard</h2>
            <div style={{ display: "flex", justifyContent: "flex-end", marginBottom: "20px" }}>
                <button onClick={logout}>Log out</button>
            </div>

            <div className="form-group">
                <input
                    placeholder="Title"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                />
                <textarea
                    placeholder="Body"
                    value={body}
                    onChange={(e) => setBody(e.target.value)}
                    rows={3}
                />
                <button onClick={addNote}>Add Note</button>
            </div>

            <ul>
                {notes.map((note) => (
                    <li key={note.id}>
                        {editNoteId === note.id ? (
                            <div style={{ flex: 1 }}>
                                <input
                                    value={editTitle}
                                    onChange={(e) => setEditTitle(e.target.value)}
                                    placeholder="Title"
                                />
                                <textarea
                                    value={editBody}
                                    onChange={(e) => setEditBody(e.target.value)}
                                    placeholder="Body"
                                    rows={2}
                                />
                            </div>
                        ) : (
                            <span>
                                <b>{note.title}</b>
                                <br />
                                {note.body}
                            </span>
                        )}

                        {editNoteId === note.id ? (
                            <button onClick={saveEdit}>
                                <FontAwesomeIcon icon={faSave} />
                            </button>
                        ) : (
                            <>
                                <button onClick={() => startEdit(note)}>
                                    <FontAwesomeIcon icon={faEdit} />
                                </button>
                                <button onClick={() => deleteNote(note.id)}>
                                    <FontAwesomeIcon icon={faTrash} />
                                </button>
                            </>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );

}
