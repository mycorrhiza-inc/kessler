from typing import Tuple
from pathlib import Path


from .niclib import rand_string

from .llm_prompts import LLM

import subprocess

import requests
import logging

from typing import List, Optional, Any

from warnings import warn

import os

import aiohttp
import asyncio

OS_TMPDIR = Path(os.environ["TMPDIR"])
DATALAB_API_KEY = os.environ["DATALAB_API_KEY"]


# Downsample audio before sending to server, human words dont convey that much information anyway
def downsample_audio(
    filepath: Path, file_type: str, bitrate: int, tmpdir: Path
) -> Path:
    outfile = tmpdir / Path(rand_string() + ".opus")
    """
    Converts an input audio or video file to a Opus audio file w/ a specified bit rate.
    """
    ffmpeg_command = [
        "ffmpeg",  # The command to invoke FFmpeg
        "-i",
        filepath,  # The input file
        "-c:a",
        "libopus",  # Codec to use for the audio conversion
        "-b:a",
        str(bitrate),  # Bitrate for the output audio
        "-vn",  # No video (discard video data)
        outfile,  # Name of the output file
    ]
    # Execute the FFmpeg command
    result = subprocess.run(
        ffmpeg_command, stdout=subprocess.PIPE, stderr=subprocess.PIPE
    )

    # Check if FFmpeg command execution was successful
    if result.returncode != 0:
        warn(
            f"Error converting video file, falling back to original. FFmpeg said:\n{result.stderr.decode()}"
        )
        return filepath
    return outfile


MARKER_ENDPOINT_URL = os.environ["MARKER_ENDPOINT_URL"]


class GPUComputeEndpoint:
    def __init__(
        self,
        logger: Any,
        marker_endpoint_url: str = MARKER_ENDPOINT_URL,
        legacy_endpoint_url: str = "https://depricated-url.com",
        datalab_api_key: str = DATALAB_API_KEY,
    ):
        self.logger = logger
        self.marker_endpoint_url = marker_endpoint_url
        self.endpoint_url = legacy_endpoint_url
        self.datalab_api_key = datalab_api_key

    async def transcribe_pdf(
        self, filepath: Path, external_process: bool = True
    ) -> str:
        if False:
            url = "https://www.datalab.to/api/v1/marker"
            self.logger.info(
                "Calling datalab api with key beginning with"
                + self.datalab_api_key[0 : (len(self.datalab_api_key) // 5)]
            )
            headers = {"X-Api-Key": self.datalab_api_key}

            async with aiohttp.ClientSession() as session:
                with aiohttp.MultipartWriter() as mpwriter:
                    with open(filepath, "rb") as file:
                        # # Append the file
                        # file_part = mpwriter.append(file)
                        # file_part.set_content_disposition(
                        #     "attachment", filename=filepath.name + ".pdf"
                        # )
                        # file_part.headers[aiohttp.hdrs.CONTENT_TYPE] = "application/pdf"

                        # # Append other form data
                        # mpwriter.append_form(
                        #     [("langs", "en"), ("force_ocr", "false"), ("paginate", "true")]
                        # )

                        # async with session.post(url, data=mpwriter, headers=headers) as response:
                        files = {
                            "file": (filepath.name + ".pdf", file, "application/pdf"),
                            "paginate": (None, True),
                        }
                        # data = {"langs": "en", "force_ocr": "false", "paginate": "true"}
                        with requests.post(
                            url, files=files, headers=headers
                        ) as response:
                            response.raise_for_status()
                            # await the json if async
                            data = response.json()
                            request_check_url = data.get("request_check_url")

                            if request_check_url is None:
                                raise Exception(
                                    "Failed to get request_check_url from marker API response"
                                )
                            self.logger.info(
                                "Got response from marker server, polling to see when file is finished processing."
                            )

                            # Polling for the result
                            max_polls = 300
                            for _ in range(max_polls):
                                await asyncio.sleep(2)
                                async with session.get(
                                    request_check_url, headers=headers
                                ) as poll_response:
                                    poll_response.raise_for_status()
                                    poll_data = await poll_response.json()

                                    if poll_data["status"] == "complete":
                                        return poll_data["markdown"]

                            raise TimeoutError(
                                "Polling for marker API result timed out"
                            )
        else:
            url = self.marker_endpoint_url + "/api/v1/marker"

            async with aiohttp.ClientSession() as session:
                with aiohttp.MultipartWriter() as mpwriter:
                    with open(filepath, "rb") as file:
                        # # Append the file
                        # file_part = mpwriter.append(file)
                        # file_part.set_content_disposition(
                        #     "attachment", filename=filepath.name + ".pdf"
                        # )
                        # file_part.headers[aiohttp.hdrs.CONTENT_TYPE] = "application/pdf"

                        # # Append other form data
                        # mpwriter.append_form(
                        #     [("langs", "en"), ("force_ocr", "false"), ("paginate", "true")]
                        # )

                        # async with session.post(url, data=mpwriter, headers=headers) as response:
                        files = {
                            "file": (filepath.name + ".pdf", file, "application/pdf"),
                            "paginate": (None, True),
                        }
                        # data = {"langs": "en", "force_ocr": "false", "paginate": "true"}
                        with requests.post(url, files=files) as response:
                            response.raise_for_status()
                            # await the json if async
                            data = response.json()
                            request_check_url_leaf = data.get("request_check_url_leaf")

                            if request_check_url_leaf is None:
                                raise Exception(
                                    "Failed to get request_check_url from marker API response"
                                )
                            request_check_url = (
                                self.marker_endpoint_url + request_check_url_leaf
                            )
                            self.logger.info(
                                "Got response from marker server, polling to see when file is finished processing."
                            )

                            # Polling for the result
                            max_polls = 300
                            for _ in range(max_polls):
                                await asyncio.sleep(2)
                                async with session.get(
                                    request_check_url
                                ) as poll_response:
                                    poll_response.raise_for_status()
                                    poll_data = await poll_response.json()

                                    if poll_data["status"] == "complete":
                                        return poll_data["markdown"]

                            raise TimeoutError(
                                "Polling for marker API result timed out"
                            )

    def llm_nonlocal_raw_call(
        self, msg_history: list[dict], model_name: Optional[str]
    ) -> dict:
        if model_name == None:
            model_name = "meta-llama-3-8b-instruct"
        if model_name == "small":
            model_name = "meta-llama-3-8b-instruct"
        if model_name == "large":
            model_name = "meta-llama-3-70b-instruct"
        # The API endpoint you will be hitting
        url = f"{self.endpoint_url}/v0/chat_completion/external_api"
        jsonpayload = {
            "messages": msg_history,
            "model_name": model_name,
        }
        print(f"Calling external endpoint: {self.endpoint_url}")
        # Make the POST request with files
        response = requests.post(url, json=jsonpayload)
        # Raise an exception if the request was unsuccessful
        response.raise_for_status()
        # Parse the JSON response
        response_json = response.json()
        # Extract the translated text from the JSON response
        return response_json["message"]

    def llm_from_model_name(self, model_name: Optional[str] = None):
        return LLM(lambda messages: self.llm_nonlocal_raw_call(messages, model_name))

    def audio_to_text_raw(
        self, filepath: Path, source_lang: str, target_lang: str, file_type: str
    ) -> str:
        # The API endpoint you will be hitting
        url = f"{self.endpoint_url}/v0/multimodal_asr/whisper-latest"
        # Open the file in binary mode
        with filepath.open("rb") as file:
            # Define the multipart/form-data payload
            files = {"file": (filepath.name, file, "application/octet-stream")}
            jsonpayload = {
                "source_lang": source_lang,
                "target_lang": target_lang,
            }
            # Make the POST request with files
            response = requests.post(url, files=files, json=jsonpayload)
            # Raise an exception if the request was unsuccessful
            response.raise_for_status()

        # Parse the JSON response
        response_json = response.json()

        # Extract the translated text from the JSON response
        translated_text = response_json["response"]
        return translated_text

    def audio_to_text(
        self,
        filepath: Path,
        source_lang: str,
        target_lang: str,
        file_type: str,
        bitrate: int = 15000,
        tmpdir: Path = OS_TMPDIR,
    ) -> str:
        downsampled = downsample_audio(filepath, file_type, bitrate, tmpdir)
        if downsampled == filepath:
            return self.audio_to_text_raw(filepath, source_lang, target_lang, file_type)
        return self.audio_to_text_raw(downsampled, source_lang, target_lang, "opus")

    def embed_raw_dicts(self, text_list: List[dict], model_name: str) -> list:
        if not model_name in ["mistral7b-sfr"]:
            raise Exception("Invalid Model ID")
        if len(text_list) == 0:
            return []
        url = f"{self.endpoint_url}/v0/embedding/{model_name}"
        payload = {"embeddable": text_list, "model_name": model_name}
        response = requests.post(url, json=payload)
        response.raise_for_status()
        embeddings = response.json().get("embeddings", [])
        print(embeddings)
        return embeddings

    def embed_queries_and_texts(
        self, query_list: List[str], text_list: List[str]
    ) -> Tuple[list, list]:
        query_dict_list = map(lambda x: {"text": x, "query": True}, query_list)
        text_dict_list = map(lambda x: {"text": x, "query": False}, text_list)

        raw_list = list(query_dict_list) + list(text_dict_list)
        embeddings = self.embed_raw_dicts(raw_list, "mistral7b-sfr")
        query_embeddings = embeddings[: len(query_list)]
        text_embeddings = embeddings[len(query_list) :]
        return (query_embeddings, text_embeddings)

    def embedding_query(self, embedding: str):
        return (self.embed_queries_and_texts([embedding], [])[0])[0]

    def translate_text(
        self, doctext: str, source_lang: Optional[str], target_lang: str
    ) -> str:
        url = f"{self.endpoint_url}/v0/translation/google-translate"
        payload = {
            "text": doctext,
            "source_lang": source_lang,
            "target_lang": target_lang,
        }
        response = requests.post(url, json=payload)
        response.raise_for_status()
        text = response.json().get("text", [])
        return text
