# Минимальный рантайм-образ
FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарник и .env (если нужно)
COPY bookingkart-platform .
COPY .env .  # <-- можно убрать, если .env уже на сервере

EXPOSE 8080

ENTRYPOINT ["./bookingkart-platform"]
