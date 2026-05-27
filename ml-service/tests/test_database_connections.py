from __future__ import annotations

from app.metrics.database import Database as MetricsDatabase
from app.ml.database import Database as MLDatabase


class FakeConnection:
    def __enter__(self) -> "FakeConnection":
        return self

    def __exit__(self, *args: object) -> None:
        return None


def test_metrics_database_disables_prepared_statements(monkeypatch) -> None:
    calls = []

    def fake_connect(*args, **kwargs):
        calls.append((args, kwargs))
        return FakeConnection()

    monkeypatch.setattr("app.metrics.database.psycopg.connect", fake_connect)

    with MetricsDatabase("postgresql://example").connect():
        pass

    assert calls[0][1]["prepare_threshold"] is None


def test_ml_database_disables_prepared_statements(monkeypatch) -> None:
    calls = []

    def fake_connect(*args, **kwargs):
        calls.append((args, kwargs))
        return FakeConnection()

    monkeypatch.setattr("app.ml.database.psycopg.connect", fake_connect)

    with MLDatabase("postgresql://example").connect():
        pass

    assert calls[0][1]["prepare_threshold"] is None
