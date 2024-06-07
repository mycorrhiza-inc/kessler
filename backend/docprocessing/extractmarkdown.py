from util.niclib import rand_string, rand_filepath
import logging

# Note: Refactoring imports.py

import requests

from typing import Optional, List, Union


# from langchain.vectorstores import FAISS

import subprocess
import os


from pathlib import Path

from util.gpu_compute_calls import GPUComputeEndpoint


import yaml

from util.niclib import create_markdown_string, seperate_markdown_string

class MarkdownExtractor:
    def __init__(self, logger, tmpdir: Path):
        self.tmpdir = tmpdir
        self.logger = logger
        # TODO : Add database connection.

    def process_raw_document_into_english_text(self, file_loc: Path, metadata: str):
        lang = metadata.get("lang")
        raw_text = self.process_raw_document_into_untranslated_text(file_loc, metadata)
        return self.convert_text_into_eng(raw_text, lang)

    def convert_text_into_eng(self, file_text: str, lang: str):
        if lang in ["en", "eng", "english", None]:
            return file_text
        english_text = GPUComputeEndpoint().translate_text(
            file_text, lang, "en"
        )
        return english_text

    def backup_processed_text(self,text : str, metadata : dict, backupdir : Path) -> None:
        savestring = create_markdown_string(text,metadata,include_previous_metadata = False)
        backuppath = backupdir / Path(metadata["hash"] + ".md")
        with open(backuppath, "w") as text_file:
            text_file.write(savestring)



    def process_raw_document_into_untranslated_text(
        self, file_loc: Path, metadata: dict,override_dir : Optional[Path] = None
    ) -> str:
        doctype = metadata["doctype"]
        lang = metadata["lang"]
            
        def process_pdf(filepath: Path) -> str:
            return GPUComputeEndpoint().transcribe_pdf(filepath)

        # Take a file with a path of path and a pandoc type of doctype and convert it to pandoc markdown and return the output as a string.
        # TODO: Make it so that you dont need to run sudo apt install pandoc for it to work, and it bundles with the pandoc python library
        def process_pandoc(filepath: Path, doctype: str) -> str:
            command = f"pandoc -f {doctype} {filepath}"
            process = subprocess.Popen(
                command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )
            output, error = process.communicate()
            output_str = output.decode()
            error_str = error.decode()
            if error_str:  # TODO : Debug this weird if statement
                raise Exception(f"Error running pandoc command: {error_str}")
            return output_str
        if not override_dir is None:
            hash = metadata["hash"]
            checkpath = override_dir / Path(hash + ".md")
            if os.path.exists(checkpath):
                with open(checkpath, "r") as file:
                    data = file.read().rstrip()
                    text,new_metadata = seperate_markdown_string(data)
                    metadata.update(new_metadata)
                    return (text,metadata)

        if not os.path.isfile(file_loc):
            raise Exception("A document with that hash is not present")
        if doctype == "md":
            with open(file_loc, "r") as file:
                data = file.read().rstrip()
                text,new_metadata = seperate_markdown_string(data)
                # Redundant due to it processing metadata upon ingest.
                # metadata.update(new_metadata)
                return (text,metadata)
            
        elif doctype == "pdf":
            return (process_pdf(file_loc),metadata)
        elif doctype in [
            "html",
            "doc",
            "docx",
            "tex",
            "epub",
            "odt",
            "rtf",
        ]:
            return (process_pandoc(file_loc, doctype),metadata)
        elif doctype == "tex":
            return (process_pandoc(file_loc, "latex"),metadata)
        elif doctype in ["mp3", "opus", "mkv"]:
            assert False
            return ""
        else:
            raise ValueError(
                f'Improper File Type, processing Failed with doctype: "{doctype}"'
            )
