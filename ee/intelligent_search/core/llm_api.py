from llama import Llama, Dialog
from decouple import config
from utils.contexts import search_context_v2
from threading import Semaphore
from asyncio import sleep
from time import time
import re


class LLM_Model:

    def __init__(self, **params):
        """
        Initialization of pre-trained model.
        Args:
            ckpt_dirckpt_dir (str): The directory containing checkpoint files for the pretrained model.
            tokenizer_path (str): The path to the tokenizer model used for text encoding/decoding.
            max_seq_len (int, optional): The maximum sequence length for input prompts. Defaults to 128.
            max_batch_size (int, optional): The maximum batch size for generating sequences. Defaults to 4.
        """
        self.generator = Llama.build(**params)
        self.max_queue_size = config('LLM_MAX_QUEUE_SIZE', cast=int, default=7)
        self.semaphore = Semaphore(config('LLM_MAX_BATCH_SIZE', cast=int, default=7))
        self.api_timeout = config('API_TIMEOUT', cast=int, default=30)
        self.queue: dict[int, str] = dict()
        self.responses: dict[int, dict] = dict()

    def __execute_prompts(self, prompts: list, **params):
        """
        Entry point of the program for generating text using a pretrained model.

        Args:
            prompts (list str): batch of prompts to be asked to LLM.
            temperature (float, optional): The temperature value for controlling randomness in generation. Defaults to 0.6.
            top_p (float, optional): The top-p sampling parameter for controlling diversity in generation. Defaults to 0.9.
            max_gen_len (int, optional): The maximum length of generated sequences. Defaults to 64.
        """
        return self.generator.text_completion(
                prompts, **params)

    def execute_prompts(self, prompts, **params):
        if self.semaphore.acquire(timeout=10):
            results = self.__execute_prompts(prompts, **params)
            self.semaphore.release()
            return results
        else:
            raise TimeoutError("[Error] LLM is over-requested")

    def __execute_prompts_dict(self, prompts: dict, **params):
        """
        Entry point of the program for generating text using a pretrained model.

        Args:
            prompts (list str): batch of prompts to be asked to LLM.
            temperature (float, optional): The temperature value for controlling randomness in generation. Defaults to 0.6.
            top_p (float, optional): The top-p sampling parameter for controlling diversity in generation. Defaults to 0.9.
            max_gen_len (int, optional): The maximum length of generated sequences. Defaults to 64.
        """
        responses = self.generator.text_completion(
                prompts.values(), **params)
        self.responses = dict(zip(prompts.keys(), responses))

    def __is_ready(self, process_id: int):
        return process_id in self.responses.keys()

    def add_to_list(self, key, value, **params):
        self.queue[key] = value
        if len(self.queue) >= self.max_queue_size:
            self.__execute_prompts_dict(self.queue, **params)
            self.queue = dict()

    def add_to_list_multiple(self, key, value, context, **params):
        self.queue[key] = value
        if len(self.queue) >= self.max_queue_size:
            results = self.execute_multiple_questions(self.queue, context, **params)
            m = re.compile('\([0-9]+\).*\n.*;')
            results = m.findall(results[0]['generation'])
            for i, k in enumerate(self.queue.keys()):
                try:
                    self.responses[k] = results[i]
                except KeyError:
                    self.responses[k] = ''
            self.queue = dict()

    async def queue_prompt(self, prompt:str, process_id: int, force=False, **params):
        if self.semaphore.acquire(timeout=40):
            if force:
                self.__execute_prompts_dict(self.queue | {process_id: prompt}, **params)
                self.queue = dict()
            else:
                self.add_to_list(process_id, prompt)
                # Wait until response exists
                init_time = time()
                while not self.__is_ready(process_id):
                    await sleep(0.1)
                    if time()-init_time > self.api_timeout:
                        raise Exception('[TimeOut] TimeoutError')
            response = self.responses.pop(process_id)
            self.semaphore.release()
            return response
        else:
            raise TimeoutError("[Error] LLM is over-requested")

    def execute_multiple_questions(self, prompts, context, **params):
        long_prompt = ''
        for i, prompt in enumerate(prompts.values()):
            long_prompt += f'({i}) {prompt}\n'
        return self.__execute_prompts([context.format(user_question=long_prompt)], **params)

    async def queue_multiple(self, prompt, context, process_id: int, force=False, **params):
        if self.semaphore.acquire(timeout=40):
            if force:
                self.execute_multiple_questions(self.queue | {process_id: prompt}, context, **params)
                self.queue = dict()
            else:
                self.add_to_list_multiple(process_id, prompt, context, **params)
                # Wait until response exists
                t = time()
                while not self.__is_ready(process_id):
                    await sleep(0.1)
                    if time()-t > self.api_timeout:
                        raise Exception('[TimeOut] TimeOutError')
            response = self.responses.pop(process_id)
            self.semaphore.release()
            return response
        else:
            raise TimeoutError("[Error] LLM is over-requested")

