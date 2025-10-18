import asyncio
import aiohttp
import json
import time
import subprocess
import os
from pathlib import Path
from typing import List, Dict, Any
import statistics

class TestRunner:
    def __init__(self, compose_file: str, test_command: str, output_dir: str, 
                 parallel_runs: int = 5, total_iterations: int = 100):
        self.compose_file = compose_file
        self.test_command = test_command
        self.output_dir = Path(output_dir)
        self.parallel_runs = parallel_runs
        self.total_iterations = total_iterations
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
    async def start_test_environment(self) -> bool:
        """Запускает тестовое окружение Docker"""
        try:
            print("Starting test environment...")
            process = await asyncio.create_subprocess_exec(
                "docker", "compose", "-f", self.compose_file, "up", "-d",
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            stdout, stderr = await process.communicate()
            
            if process.returncode == 0:
                print("Test environment started successfully")
                # Даем время для запуска сервисов
                await asyncio.sleep(10)
                return True
            else:
                print(f"Failed to start test environment: {stderr.decode()}")
                return False
        except Exception as e:
            print(f"Error starting test environment: {e}")
            return False
    
    async def stop_test_environment(self):
        """Останавливает и удаляет контейнеры тестового окружения"""
        try:
            print("Stopping test environment...")
            process = await asyncio.create_subprocess_exec(
                "docker", "compose", "-f", self.compose_file, "down", '--remove-orphans',
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            await process.communicate()
            print("Test environment stopped")
        except Exception as e:
            print(f"Error stopping test environment: {e}")
    
    async def run_single_test(self, test_id: int) -> Dict[str, Any]:
        """Запускает один тест и возвращает время выполнения"""
        start_time = time.time()
        try:
            process = await asyncio.create_subprocess_exec(
                *self.test_command.split(),
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            stdout, stderr = await process.communicate()
            end_time = time.time()
            
            execution_time = end_time - start_time
            success = process.returncode == 0
            
            return {
                "test_id": test_id,
                "execution_time": execution_time,
                "success": success,
                "stdout": stdout.decode() if stdout else "",
                "stderr": stderr.decode() if stderr else ""
            }
        except Exception as e:
            return {
                "test_id": test_id,
                "execution_time": time.time() - start_time,
                "success": False,
                "error": str(e)
            }
    
    async def get_metrics(self) -> Dict[str, Any]:
        """Запрашивает метрики с локального сервера"""
        url = "http://localhost:23450/api/metrics/get"
        try:
            async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=10)) as session:
                async with session.get(url) as response:
                    if response.status == 200:
                        return await response.json()
                    else:
                        print(f"Failed to get metrics: HTTP {response.status}")
                        return {}
        except Exception as e:
            print(f"Error getting metrics: {e}")
            return {}
    
    async def run_parallel_tests(self, iteration: int) -> Dict[str, Any]:
        """Запускает параллельные тесты для одной итерации"""
        print(f"Iteration {iteration}: Running {self.parallel_runs} parallel tests...")
        
        # Запускаем все тесты параллельно
        tasks = [self.run_single_test(i + 1) for i in range(self.parallel_runs)]
        test_results = await asyncio.gather(*tasks)
        
        # Собираем времена выполнения
        execution_times = [result["execution_time"] for result in test_results]
        successful_tests = sum(1 for result in test_results if result["success"])
        
        # Получаем метрики
        metrics = await self.get_metrics()
        
        return {
            "iteration": iteration,
            "timestamp": time.time(),
            "test_results": {
                "execution_times": execution_times,
                "successful_tests": successful_tests,
                "total_tests": len(test_results),
                "min_time": min(execution_times) if execution_times else 0,
                "max_time": max(execution_times) if execution_times else 0,
                "avg_time": statistics.mean(execution_times) if execution_times else 0
            },
            "metrics": metrics.get("metrics", {}) if metrics else {}
        }
    
    def save_iteration_results(self, iteration: int, results: Dict[str, Any]):
        """Сохраняет результаты итерации в JSON файл"""
        filename = self.output_dir / f"run_{iteration}.json"
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        print(f"Results saved to {filename}")
    
    async def run(self):
        """Основной метод запуска тестов"""
        print(f"Starting test run: {self.total_iterations} iterations, {self.parallel_runs} parallel tests each")
        
        all_results = []
        
        for iteration in range(1, self.total_iterations + 1):
            try:
                # Запускаем тестовое окружение
                if not await self.start_test_environment():
                    print(f"Failed to start test environment for iteration {iteration}")
                    continue
                
                # Запускаем параллельные тесты и собираем результаты
                iteration_results = await self.run_parallel_tests(iteration)
                all_results.append(iteration_results)
                
                # Сохраняем результаты итерации
                self.save_iteration_results(iteration, iteration_results)
                
                # Останавливаем тестовое окружение
                await self.stop_test_environment()
                
                # Небольшая пауза между итерациями
                await asyncio.sleep(2)
                
            except Exception as e:
                print(f"Error in iteration {iteration}: {e}")
                await self.stop_test_environment()
                continue
        
        # Сохраняем сводную статистику
        self.save_summary_statistics(all_results)
        
        print("Test run completed!")
    
    def save_summary_statistics(self, all_results: List[Dict[str, Any]]):
        """Сохраняет сводную статистику по всем итерациям"""
        if not all_results:
            print("No results to summarize")
            return
        
        # Собираем данные о времени выполнения тестов
        all_avg_times = [result["test_results"]["avg_time"] for result in all_results]
        all_min_times = [result["test_results"]["min_time"] for result in all_results]
        all_max_times = [result["test_results"]["max_time"] for result in all_results]
        
        # Собираем метрики CPU и Memory
        cpu_current_values = [result["metrics"].get("cpu_current_value", 0) for result in all_results]
        cpu_max_values = [result["metrics"].get("cpu_max_value", 0) for result in all_results]
        cpu_min_values = [result["metrics"].get("cpu_min_value", 0) for result in all_results]
        cpu_avg_values = [result["metrics"].get("cpu_avg_value", 0) for result in all_results]
        
        memory_current_values = [result["metrics"].get("memory_current_value", 0) for result in all_results]
        memory_max_values = [result["metrics"].get("memory_max_value", 0) for result in all_results]
        memory_min_values = [result["metrics"].get("memory_min_value", 0) for result in all_results]
        memory_avg_values = [result["metrics"].get("memory_avg_value", 0) for result in all_results]
        
        summary = {
            "total_iterations": len(all_results),
            "test_execution_time_summary": {
                "min_execution_time": min(all_min_times) if all_min_times else 0,
                "max_execution_time": max(all_max_times) if all_max_times else 0,
                "avg_execution_time": statistics.mean(all_avg_times) if all_avg_times else 0,
                "median_execution_time": statistics.median(all_avg_times) if all_avg_times else 0
            },
            "cpu_metrics_summary": {
                "cpu_current_avg": statistics.mean(cpu_current_values) if cpu_current_values else 0,
                "cpu_max_avg": statistics.mean(cpu_max_values) if cpu_max_values else 0,
                "cpu_min_avg": statistics.mean(cpu_min_values) if cpu_min_values else 0,
                "cpu_avg_avg": statistics.mean(cpu_avg_values) if cpu_avg_values else 0
            },
            "memory_metrics_summary": {
                "memory_current_avg": statistics.mean(memory_current_values) if memory_current_values else 0,
                "memory_max_avg": statistics.mean(memory_max_values) if memory_max_values else 0,
                "memory_min_avg": statistics.mean(memory_min_values) if memory_min_values else 0,
                "memory_avg_avg": statistics.mean(memory_avg_values) if memory_avg_values else 0
            },
            "iterations_completed": len(all_results)
        }
        
        summary_file = self.output_dir / "summary_statistics.json"
        with open(summary_file, 'w', encoding='utf-8') as f:
            json.dump(summary, f, indent=2, ensure_ascii=False)
        
        print(f"Summary statistics saved to {summary_file}")
        print(f"Completed {len(all_results)} out of {self.total_iterations} iterations")

async def main():
    # Конфигурация
    config = {
        "compose_file": "./deployments/docker-compose.test.yaml",
        "test_command": "make test-e2e",
        "output_dir": "./test_results/debug_with_tracing",
        "parallel_runs": 5, 
        "total_iterations": 100 
    }
    
    # Создаем и запускаем тестер
    runner = TestRunner(
        compose_file=config["compose_file"],
        test_command=config["test_command"],
        output_dir=config["output_dir"],
        parallel_runs=config["parallel_runs"],
        total_iterations=config["total_iterations"]
    )
    
    await runner.run()

if __name__ == "__main__":
    asyncio.run(main())
