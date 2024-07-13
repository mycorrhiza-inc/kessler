from rag.llamaindex import add_document_to_db_from_text
import os
from pathlib import Path
from typing import Any

from sqlalchemy import select

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from models.files import (
    FileModel,
    FileRepository,
    FileSchema,
    FileSchemaWithText,
    provide_files_repo,
    DocumentStatus,
    docstatus_index
)

from logic.documentingest import DocumentIngester
from logic.extractmarkdown import MarkdownExtractor

from typing import List, Optional, Dict


import json

from util.niclib import rand_string


async def add_file_raw(
    tmp_filepath: Path,
    metadata: dict,
    process: bool ,  # Figure out how to pass in a boolean as a query paramater
    override_hash: bool,
    files_repo: FileRepository,
    logger: Any,
) -> FileSchema:
    docingest = DocumentIngester(logger)

    def validate_metadata_mutable(metadata: dict):
        if metadata.get("lang") is None:
            metadata["lang"] = "en"
        try:
            assert isinstance(metadata.get("title"), str)
            assert isinstance(metadata.get("doctype"), str)
            assert isinstance(metadata.get("lang"), str)
        except Exception:
            logger.error("Illformed Metadata please fix")
            logger.error(f"Title: {metadata.get("title")}")
            logger.error(f"Doctype: {metadata.get("doctype")}")
            logger.error(f"Lang: {metadata.get("title")}")
            raise Exception(
                "Metadata is illformed, this is likely an error in software, please submit a bug report."
            )
        else:
            logger.info("Title, Doctype and language successfully declared")

        if (metadata["doctype"])[0] == ".":
            metadata["doctype"] = (metadata["doctype"])[1:]
        if metadata.get("source") is None:
            metadata["source"] = "unknown"
        metadata["language"] = metadata["lang"]
        return metadata

    # This assignment shouldnt be necessary, but I hate mutating variable bugs.
    metadata = validate_metadata_mutable(metadata)

    logger.info("Attempting to save data to file")
    result = docingest.save_filepath_to_hash(tmp_filepath, OS_HASH_FILEDIR)
    (filehash, filepath) = result

    os.remove(tmp_filepath)

    # NOTE: this is a dangeous query
    # NOTE: Nicole- Also this doesnt allow for files with the same hash to have different metadata,
    # Scrapping it is a good idea, it was designed to solve some issues I had during testing and adding the same dataset over and over
    # FIX: fix this to not allow for users to DOS files
    query = select(FileModel).where(FileModel.hash == filehash)

    duplicate_file_objects = await files_repo.session.execute(query)
    duplicate_file_obj = duplicate_file_objects.scalar()

    if override_hash is True and duplicate_file_obj is not None:
        try:
            await files_repo.delete(duplicate_file_obj.id)
        except Exception:
            pass
        duplicate_file_obj = None

    if duplicate_file_obj is None:
        docingest.backup_metadata_to_hash(metadata, filehash)
        metadata_str = json.dumps(metadata)
        new_file = FileModel(
            url="N/A",
            name=metadata.get("title"),
            doctype=metadata.get("doctype"),
            lang=metadata.get("lang"),
            source=metadata.get("source"),
            mdata=metadata_str,
            stage=(DocumentStatus.stage1).value,
            hash=filehash,
            summary=None,
            short_summary=None,
        )
        logger.info("new file:{file}".format(file=new_file.to_dict()))
        try:
            new_file = await files_repo.add(new_file)
        except Exception as e:
            logger.info(e)
            return e
        logger.info("added file!~")
        await files_repo.session.commit()
        logger.info("commited file to DB")

    else:
        logger.info(type(duplicate_file_obj))
        logger.info(
            f"File with identical hash already exists in DB with uuid:\
            {duplicate_file_obj.id}"
        )
        new_file = duplicate_file_obj

    if process:
        logger.info("Processing File")
        # self.process_file_raw(new_file, files_repo, logger, "")
        # Removing the await here, SHOULD allow an instant return to the user, but also have python process the file in the background hopefully!
        # TODO : It doesnt work, for now its fine and also doesnt compete with the daemon for resources. Fix later
        await process_file_raw(new_file, files_repo, logger, "")

    return new_file


async def process_file_raw(
    self, obj: FileModel, files_repo: FileRepository, logger: Any, regenerate: Optional[str] = None, stop_at : Optional[str]= None
):
    if stop_at is None:
        stop_at = "completed"
    if regenerate is None:
        regenerate = "completed"
    # TODO: Figure out error messaging for these
    stop_at=DocumentStatus(stop_at)
    regenerate=DocumentStatus(regenerate)


    logger.info(type(obj))
    logger.info(obj)
    current_stage = obj.stage
    logger.info(obj.doctype)
    mdextract = MarkdownExtractor(logger, OS_TMPDIR)
    doc_metadata = json.loads(obj.mdata)

    response_code, response_message = (
        500,
        "Internal error somewhere in process.",
    )

    if docstatus_index(current_stage) < docstatus_index(regenerate):
        current_stage = regenerate

    # TODO: Replace with pydantic validation

    # text extraction
    def process_stage_one():
        # FIXME: Change to deriving the filepath from the uri.
        file_path = DocumentIngester(logger).get_default_filepath_from_hash(
            obj.hash
        )
        # This process might spit out new metadata that was embedded in the document, ignoring for now
        processed_original_text = (
            mdextract.process_raw_document_into_untranslated_text(
                file_path, doc_metadata
            )[0]
        )
        logger.info(
            f"Successfully processed original text: {processed_original_text[0:20]}"
        )
        # FIXME: We should probably come up with a better backup protocol then doing everything with hashes
        mdextract.backup_processed_text(
            processed_original_text, obj.hash, doc_metadata, OS_BACKUP_FILEDIR
        )
        assert isinstance(processed_original_text, str)
        logger.info("Backed up markdown text")
        if obj.lang == "en":
            # Write directly to the english text box if
            # original text is identical to save space.
            obj.english_text = processed_original_text
            # Skip translation stage if text already english.
            return DocumentStatus.stage3
        else:
            obj.original_text = processed_original_text
            return DocumentStatus.stage2

    # text conversion
    def process_stage_two():
        if obj.lang != "en":
            try:
                processed_english_text = mdextract.convert_text_into_eng(
                    obj.original_text, obj.lang
                )
                obj.english_text = processed_english_text
            except Exception as e:
                raise Exception(
                    "failure in stage 2: \ndocument was unable to be translated to english.",
                    e,
                )
        else:
            raise ValueError(
                "failure in stage 2: \n Code is in an unreachable state, a document cannot be english and not english",
            )
        return DocumentStatus.stage3

    # TODO: Replace with pydantic validation

    def process_stage_three():
        logger.info("Adding Document to Vector Database")

        def generate_searchable_metadata(initial_metadata : dict) -> dict:
            return_metadata = {
                "title"  : initial_metadata.get("title"),
                "author" : initial_metadata.get("author"),
                "source" : initial_metadata.get("source"),
                "date" : initial_metadata.get("date"),
            }
            def guarentee_field(field: str, default_value : Any = "unknown"):
                if return_metadata.get(field) is None:
                    return_metadata[field]=default_value
            guarentee_field("title")
            guarentee_field("author")
            guarentee_field("source")
            guarentee_field("date")
            return return_metadata
        searchable_metadata = generate_searchable_metadata(doc_metadata)
        try:
            add_document_to_db_from_text(obj.english_text, searchable_metadata)
        except Exception as e:
            raise Exception("Failure in adding document to vector database", e)
        return DocumentStatus.completed

    while True:
        if docstatus_index(current_stage) <= docstatus_index(stop_at):
            response_code, response_message = (
                200,
                "Document Fully Processed.",
            )
            logger.info(current_stage.value)
            obj.stage = current_stage.value
            logger.info(response_code)
            logger.info(response_message)
            _ = files_repo.update(obj)
            await files_repo.session.commit()
            type_adapter = TypeAdapter(FileSchema)
            final_return = type_adapter.validate_python(obj)
            return final_return

        match current_stage:
            case DocumentStatus.stage1:
                current_stage = process_stage_one()
            case DocumentStatus.stage2:
                current_stage = process_stage_two()
            case DocumentStatus.stage3:
                current_stage = process_stage_three()
            case _:
                raise Exception(
                    "Document was incorrectly added to database, \
                    try readding it again.\
                "
                )
            # FIXME: The try catch exception code broke the plaintext error handling, since it still returns a 500 error, I removed it temporarially
            # try:
            #    match current_stage:
            #        case "stage1":
            #            current_stage = process_stage_one()
            #        case "stage2":
            #            current_stage = process_stage_two()
            #        case "stage3":
            #            current_stage = process_stage_three()
            #        case "completed":
            #            response_code, response_message = (
            #                200,
            #                "Document Fully Processed.",
            #            )
            #            logger.info(current_stage)
            #            obj.stage = current_stage
            #            logger.info(response_code)
            #            logger.info(response_message)
            #            _ = files_repo.update(obj)
            #            await files_repo.session.commit()
            #            break
            #        case _:
            #            raise Exception(
            #                "Document was incorrectly added to database, \
            #                try readding it again.\
            #            "
            #            )

            # except Exception as e:
            #    logger.error(e)
            #    break