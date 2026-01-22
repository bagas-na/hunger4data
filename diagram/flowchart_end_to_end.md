```mermaid
flowchart LR

    %% =========================
    %% Client & API Gateway
    %% =========================
    User[User / Client]

    subgraph APIGW[API Gateway]
        GW_HTTP[transport/http<br/>router.go<br/>handler.go]
        GW_MW{middleware.go<br/>JWT validation}
        GW_SVC[service/gateway.go]
        GW_GRPC[adapter/grpc<br/>auth.go<br/>client.go]
    end

    User -->|REST| GW_HTTP
    GW_HTTP --> GW_MW
    GW_MW -->|valid JWT| GW_SVC
    GW_MW -->|invalid JWT| User
    GW_SVC --> GW_GRPC

    %% =========================
    %% Authentication Service
    %% =========================
    subgraph AUTH[Authentication Service]
        subgraph AUTH_1[container]
            AUTH_GRPC[transport/grpc<br/>handlers.go]
            AUTH_SVC[service/auth.go<br/>register & login]
            AUTH_REPO[adapters/repository<br/>users_postgres.go]
            AUTH_MQ[adapter/rabbitmq<br/>publisher]
        end
        DB_AUTH[(Auth DB)]
    end

    GW_GRPC -->|Login / Register| AUTH_GRPC
    AUTH_GRPC --> AUTH_SVC
    AUTH_SVC -->|JWT/Error 401| GW_GRPC

    AUTH_SVC --> AUTH_REPO
    AUTH_REPO --> DB_AUTH

    AUTH_SVC -->|Publish email_confirmation| AUTH_MQ

    %% =========================
    %% Payment Service
    %% =========================
    subgraph PAY[Payment Service]
        subgraph PAY_1[container]
            PAY_GRPC[transport/grpc<br/>handler.go]
            PAY_SVC[service/payment.go]
            PAY_MQ[adapter/rabbitmq<br/>publisher]
            PAY_REPO[adapters/repo<br/>payments_postgres.go]
            PAY_EXT[adapters/external<br/>xendit.go]
        end
        DB_PAY[(Payment DB)]
        PAYAPI[[3rd-Party Payment API]]
    end

    GW_GRPC -->|CreatePayment| PAY_GRPC
    PAY_GRPC --> PAY_SVC
    PAY_SVC --> GW_GRPC

    PAY_SVC --> PAY_EXT
    PAY_EXT --> PAYAPI

    PAY_SVC --> PAY_REPO
    PAY_REPO --> DB_PAY

    PAY_SVC -->|Publish payment_result| PAY_MQ

    %% =========================
    %% Subscription Service
    %% =========================
    subgraph SUB[Subscription Service]
        subgraph SUB_1[container]
            SUB_GRPC[transport/grpc<br/>handler.go]
            SUB_SVC_SUBS[service/subscription.go]
            SUB_SVC_LIST[service/list.go]
            SUB_SCHED[service/scheduler.go]
            SUB_REPO[adapters/repository<br/>subscriptions_postgres.go]
            SUB_EXT[adapters/external<br/>humdata.go]
            SUB_MQ[adapters/rabbitmq<br/>publisher.go]
            SUB_CRON(Daily Cron / Ticker)
        end
        HUM[[hapi.humdata.org]]
        DB_SUB[(Subscription DB)]
    end

    GW_GRPC -->|Subscribe / Unsubscribe| SUB_GRPC
    SUB_GRPC --> SUB_SVC_LIST
    SUB_SVC_LIST --> GW_GRPC
    SUB_SVC_LIST --> SUB_EXT
    SUB_EXT --> HUM

    SUB_GRPC --> SUB_SVC_SUBS
    SUB_SVC_SUBS --> SUB_REPO
    SUB_REPO --> DB_SUB

    SUB_CRON --> SUB_SCHED
    SUB_SCHED --> SUB_REPO
    SUB_SCHED -->|Publish subscription_update| SUB_MQ


    %% =========================
    %% Rabbit MQ Service
    %% =========================
    RABBIT_MQ@{ shape: procs, label: "rabbitMQ"}

    AUTH_MQ --> RABBIT_MQ
    PAY_MQ --> RABBIT_MQ
    SUB_MQ --> RABBIT_MQ
    RABBIT_MQ --> NOTIF_MQ

    %% =========================
    %% Notification Service
    %% =========================
    subgraph NOTIF[Notification Service]
        subgraph NOTIF_1[container]
            NOTIF_MQ[transport/rabbitmq<br/>consumer.go]
            NOTIF_SVC[service/notification.go]
            NOTIF_REPO[adapters/repository<br/>notification_postgres.go]
            NOTIF_SMTP[adapters/smtp<br/>mailer.go]
        end
        DB_NOTIF[(Notification DB)]
        EMAIL[SMTP Provider]
    end

    NOTIF_MQ --> NOTIF_SVC

    NOTIF_SVC --> NOTIF_REPO
    NOTIF_REPO --> DB_NOTIF

    NOTIF_SVC --> NOTIF_SMTP
    NOTIF_SMTP --> EMAIL
```

[![](https://mermaid.ink/img/pako:eNqlWPtv4kYQ_lesle7ESSQ8kwBqT00ITahCwvFo2hqEFnsDq_Pr1nYSGvjfb3bttdfEpkePH0LGM9_sPL4Z27whwzUJ6qAny30x1pgF2t1o5swcDT4fPmi_Fn0Si65FiRNoH7XLYV-7wQF5wZsfhU99wnT-R6vEfubybD9crhj21tztzaOuOJ9HBvxz87i4nUyGesCw43suCyrrIPB-WbLKZ-aGAWGnK1dIa-yYlhCz6MHjm01NUL1gRqTxH48T7Rlb1MQBdZ1dBjD-s6tDvM_UIJVVFM87pzejYVfHJvYggMqKeYbwisNgLU8wRK4pkDimzFtU4-Tk83bUG0-2MsNIFwtcHQWfXB48CoyIWoP4t3Gs-wbUUUz4UYkBGEu_PP6jKHAJqUE-1BAF08ZRfX7UQdrp6eRWz_elFDhjvqjphusEmDqEKTb8I9SiFSk7kmbEfPCzvUtgapfVvjGyoj50FdhuuSvq5GFHveGD7L5fYcRzfRq4bCMchODVX8ClYMVI0emDLwl7GF4uaWB_E2AvXFrUX2cyFcyRwvXVQtSwxIsI0qf39IrbK8hwx1OAyRvFWW3TmkXGiSiYISuj6GLSbIFOlR5jLtOa1dp2n0OqaVojRcdFoYwTyANuh1H2GrExtRbQ9SfKbEGSrSzbUaQd4o3NF9f_Zevw8m99z0ceS8HsAEm59j84-p4kHKQy1IuiyDc8hkwS857CESQ65xB_Obr31yQFk1f4crAlHLwCEelenPsE5mUtybqmFI6dw21A1xvMPBnCnWqT9BAuzw9zvcsIbOvYfpvUPTKUkuBgXN1UU7AZVVWcdqoAQSogtiJIOgZSklMAcg4omYG4EwvoQWjF6RxJ_3G49A1GvZ_a2OPplZ7nKG8QwPbAIHDt0YPAQVCaBXyPk2nwlXiKIXf98SSBQE2DAtPube86dW2siRkWhnJw86thHZog7qh4gtahDU8lOB-WjrpfMOv5uO7o4b50DVt1o3WZy-8IE2p8JexTwZTeTge6vsYePZXhuGw1n2fGmDOjlKHG3u2I__N-TGPEkkAQU8eX0jbhRwSRkhgXtaWpWl7Jjm6uNi55qpXTC4mqgeaeyrmX9cuvJCbpiEtJjjjIcmBlE1LHnHaKVy4WuEx0yXJQmbYIPWhQXL9oQxyxI0aCRNrgy7ELYnR5ddWfwIG_vWk-MIV0NI-5hl_WLLwkVkeboYigEBLaZW73cBhPNPGQLsF8TZRYniYRhfL-YdL__dglee8G9OlnH2vFyXqeq7w1GcVZvChlHsqizMw64PzQzhv1CKk-ODhKTEX2B5ea6uDATotPHsCbWuLJt-MXNf44tx_u_kNBVMJSpobZZ4Pe4LJ_p_MjtCFzn6mZFE55IJC1Uwgh7vSqVt6f0-xVrTrA4sohMI8mox3E724i2JmDymjFqIk6AQtJGUHToBYgojeOmiF4EbLJDPFpMTH7OkMzZwcYDzv_uK4tYfCiu1qjzhO2fJCieb-mGNiUmkANCOu6oROgTq3eblSFF9R5Q6-oc1Jrthqn1Va10W5dNM6arfZ5GW3gerPZOG3X6rVWs9Gsts_rFxe7MvpXHN04rZ-Ddb1dbbTqjdZZ_ayMiMmJMYh-ThC_Kuy-AwBh3IM?type=png)](https://mermaid.live/edit#pako:eNqlWPtv4kYQ_lesle7ESSQ8kwBqT00ITahCwvFo2hqEFnsDq_Pr1nYSGvjfb3bttdfEpkePH0LGM9_sPL4Z27whwzUJ6qAny30x1pgF2t1o5swcDT4fPmi_Fn0Si65FiRNoH7XLYV-7wQF5wZsfhU99wnT-R6vEfubybD9crhj21tztzaOuOJ9HBvxz87i4nUyGesCw43suCyrrIPB-WbLKZ-aGAWGnK1dIa-yYlhCz6MHjm01NUL1gRqTxH48T7Rlb1MQBdZ1dBjD-s6tDvM_UIJVVFM87pzejYVfHJvYggMqKeYbwisNgLU8wRK4pkDimzFtU4-Tk83bUG0-2MsNIFwtcHQWfXB48CoyIWoP4t3Gs-wbUUUz4UYkBGEu_PP6jKHAJqUE-1BAF08ZRfX7UQdrp6eRWz_elFDhjvqjphusEmDqEKTb8I9SiFSk7kmbEfPCzvUtgapfVvjGyoj50FdhuuSvq5GFHveGD7L5fYcRzfRq4bCMchODVX8ClYMVI0emDLwl7GF4uaWB_E2AvXFrUX2cyFcyRwvXVQtSwxIsI0qf39IrbK8hwx1OAyRvFWW3TmkXGiSiYISuj6GLSbIFOlR5jLtOa1dp2n0OqaVojRcdFoYwTyANuh1H2GrExtRbQ9SfKbEGSrSzbUaQd4o3NF9f_Zevw8m99z0ceS8HsAEm59j84-p4kHKQy1IuiyDc8hkwS857CESQ65xB_Obr31yQFk1f4crAlHLwCEelenPsE5mUtybqmFI6dw21A1xvMPBnCnWqT9BAuzw9zvcsIbOvYfpvUPTKUkuBgXN1UU7AZVVWcdqoAQSogtiJIOgZSklMAcg4omYG4EwvoQWjF6RxJ_3G49A1GvZ_a2OPplZ7nKG8QwPbAIHDt0YPAQVCaBXyPk2nwlXiKIXf98SSBQE2DAtPube86dW2siRkWhnJw86thHZog7qh4gtahDU8lOB-WjrpfMOv5uO7o4b50DVt1o3WZy-8IE2p8JexTwZTeTge6vsYePZXhuGw1n2fGmDOjlKHG3u2I__N-TGPEkkAQU8eX0jbhRwSRkhgXtaWpWl7Jjm6uNi55qpXTC4mqgeaeyrmX9cuvJCbpiEtJjjjIcmBlE1LHnHaKVy4WuEx0yXJQmbYIPWhQXL9oQxyxI0aCRNrgy7ELYnR5ddWfwIG_vWk-MIV0NI-5hl_WLLwkVkeboYigEBLaZW73cBhPNPGQLsF8TZRYniYRhfL-YdL__dglee8G9OlnH2vFyXqeq7w1GcVZvChlHsqizMw64PzQzhv1CKk-ODhKTEX2B5ea6uDATotPHsCbWuLJt-MXNf44tx_u_kNBVMJSpobZZ4Pe4LJ_p_MjtCFzn6mZFE55IJC1Uwgh7vSqVt6f0-xVrTrA4sohMI8mox3E724i2JmDymjFqIk6AQtJGUHToBYgojeOmiF4EbLJDPFpMTH7OkMzZwcYDzv_uK4tYfCiu1qjzhO2fJCieb-mGNiUmkANCOu6oROgTq3eblSFF9R5Q6-oc1Jrthqn1Va10W5dNM6arfZ5GW3gerPZOG3X6rVWs9Gsts_rFxe7MvpXHN04rZ-Ddb1dbbTqjdZZ_ayMiMmJMYh-ThC_Kuy-AwBh3IM)
