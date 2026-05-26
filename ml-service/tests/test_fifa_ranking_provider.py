from __future__ import annotations

from datetime import date

from app.metrics.fifa_ranking_provider import InMemoryFifaRankingProvider


def test_fifa_ranking_provider_never_uses_future_ranking() -> None:
    provider = InMemoryFifaRankingProvider(
        [
            {"team_id": "br", "ranking_date": date(2024, 1, 1), "rank": 5},
            {"team_id": "br", "ranking_date": date(2024, 2, 1), "rank": 1},
        ]
    )

    assert provider.get_latest_ranking_before("br", date(2024, 1, 20)) == 5
    assert provider.get_latest_ranking_before("br", date(2024, 1, 1)) is None
